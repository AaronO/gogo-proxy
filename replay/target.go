package replay

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Target is a HTTP request blueprint
type Target struct {
	Method string
	URL    string
	Body   []byte
	Header http.Header
}

func NewTarget(req *http.Request) (*Target, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return &Target{
		Method: req.Method,
		Body:   data,
		Header: req.Header,
		URL:    req.URL.String(),
	}, nil
}

// Request creates an *http.Request out of Target and returns it along with an
// error in case of failure.
func (t *Target) Request() (*http.Request, error) {
	req, err := http.NewRequest(t.Method, t.URL, bytes.NewBuffer(t.Body))
	if err != nil {
		return nil, err
	}
	for k, vs := range t.Header {
		req.Header[k] = make([]string, len(vs))
		copy(req.Header[k], vs)
	}
	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	return req, nil
}
