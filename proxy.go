package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
)

type ProxyOptions struct {
	// Number of times a request should be tried
	Retries int

	// Period to wait between retries
	Period time.Duration

	// Returns a url that we should proxy to for a given request
	Balancer func(req *http.Request) (string, error)

	// A static backend to route to
	Backend string

	// Rewriter is a function modifying the request before it's proxied
	// usually it should change the hostname of the request object, etc ...
	Rewriter func(req *http.Request, destUrl *url.URL)
}

type Proxy struct {
	*ProxyOptions

	// Http proxy
	httpProxy http.Handler

	// Websocket proxy
	websocketProxy http.Handler
}

// New returns a new Proxy instance based on the provided ProxyOptions
// either 'Backend' (static) or 'Balancer' must be provided
func New(opts ProxyOptions) (*Proxy, error) {
	// Validate Balancer and Backend options
	if opts.Balancer == nil {
		if opts.Backend == "" {
			return nil, errors.New("Please provide a Backend or a Balancer")
		} else if err := validateUrl(normalizeUrl(opts.Backend)); err != nil {
			return nil, err
		} else {
			// Normalize backend's url
			opts.Backend = normalizeUrl(opts.Backend)
		}
	}

	// Default for Retries
	if opts.Retries == 0 {
		opts.Retries = 1
	}

	// Default for Period
	if opts.Period == 0 {
		opts.Period = 100 * time.Millisecond
	}

	// Default for Rewriter
	if opts.Rewriter == nil {
		opts.Rewriter = DefaultRewriter
	}

	p := &Proxy{
		ProxyOptions: &opts,
	}

	return p.init(), nil
}

// ServeHTTP allows us to comply to the http.Handler interface
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if isWebsocket(req) {
		p.websocketProxy.ServeHTTP(rw, req)
		return
	}

	p.httpProxy.ServeHTTP(rw, req)
}

// init sets up proxies and other stuff based on options
func (p *Proxy) init() *Proxy {
	// Setup http proxy
	p.httpProxy = &httputil.ReverseProxy{
		Director: p.director,
	}

	// Setup websocket proxy
	p.websocketProxy = &websocketproxy.WebsocketProxy{
		Backend: func(req *http.Request) *url.URL {
			url, _ := p.backend(req)
			p.Rewriter(req, url)
			return url
		},
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	return p
}

// director rewrites a http.Request to route to the correct host
func (p *Proxy) director(req *http.Request) {
	url, err := p.backend(req)
	if url == nil || err != nil {
		return
	}

	p.Rewriter(req, url)
}

// backend return the full url object of the backend to route
// after fetching the host backend and resolving it relative
// to the current http request data
func (p *Proxy) backend(req *http.Request) (*url.URL, error) {
	hostUrl, err := p.backendHost(req)
	if err != nil {
		return nil, err
	}

	return fullBackendURL(req, hostUrl)
}

// backendHost parses the result of getBackend and ensures it's validity
func (p *Proxy) backendHost(req *http.Request) (*url.URL, error) {
	rawurl, err := p.getBackend(req)
	if err != nil {
		return nil, err
	}

	// Normalize URL
	backendUrl := normalizeUrl(rawurl)

	if err := validateUrl(backendUrl); err != nil {
		return nil, err
	}

	return url.Parse(backendUrl)
}

// getBackend gets the backend selected by the balancer or the static one set by the 'Backend' attribute
func (p *Proxy) getBackend(req *http.Request) (string, error) {
	if p.Balancer == nil && p.Backend != "" {
		return p.Backend, nil
	}
	return p.Balancer(req)
}

// fullBackendURL builds a url.URL object representing the destination's full URL
// resolving path's etc ... relative to the request
func fullBackendURL(req *http.Request, hostUrl *url.URL) (*url.URL, error) {
	// Pick scheme differently depending on if the request is http or websocket
	schemeFunc := httpScheme
	if isWebsocket(req) {
		schemeFunc = websocketScheme
	}

	newUrl, err := url.Parse(hostUrl.String())
	if err != nil {
		return nil, err
	}

	newUrl.Scheme = schemeFunc(hostUrl.Scheme)
	newUrl.Host = hostUrl.Host
	newUrl.Path = path.Join(hostUrl.Path, req.URL.Path)

	return newUrl, nil
}

// DefaultRewriter adjust the *http.Request object it's attributes
// match those of the destionation url (scheme, host, path, ...)
func DefaultRewriter(req *http.Request, url *url.URL) {
	req.URL.Scheme = url.Scheme
	req.URL.Host = url.Host
	req.URL.Path = url.Path

	// Rewrite outgoing request url
	req.Host = url.Host
}
