package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
)

// Replayer provides us with a mechanism to replay HTTP requests
type Replayer struct {
	req *http.Request

	writer http.ResponseWriter

	target *Target

	Errors chan error

	// Current play record
	Play Play

	// List of previous plays for request
	plays []Play
}

// Represents one request try
type Play struct {
	Bytes int
	Writes int
	Status int
	Time time.Duration
}

func NewReplayer(req *http.Request, w http.ResponseWriter) (*Replayer, error) {
	r := Replayer{
		req: req,
		writer: w,
	}

	return r.init()
}

// Records request to replay it
func (r *Replayer) init() (*Replayer, error) {
	target, err := NewTarget(r.req)
	if err != nil {
		return nil, err
	}

	r.target = target

	return r, nil
}

// Reset the Replayer's state so it's ready to be replayed
func (r *Replayer) Stop() {
	r.plays = append(r.plays, r.Play)
	r.Play = Play{}
}

func (r *Replayer) Replay() (*http.Request, error) {
	return r.target.Request()
}

////
// http.ResponseWriter methods
////

func (r *Replayer) Header() http.Header {
	fmt.Printf("Getting headers")
	return r.writer.Header()
}

func (r *Replayer) Write(data []byte) (int, error) {
	fmt.Printf("Writing %d bytes\n", len(data))
	r.Play.Bytes += len(data)
	r.Play.Writes += 1
	return r.writer.Write(data)
}

func (r *Replayer) WriteHeader(status int) {
	fmt.Printf("Writing status %d\n", status)
	r.Play.Status = status
	r.writer.WriteHeader(status)
}

func (r *Replayer) GetError() error {
	if r.Play.Status >= 500 {
		return fmt.Errorf("Status code is %d", r.Play.Status)
	}

	return nil
}

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
		Body: data,
		Header: req.Header,
		URL: req.URL.String(),
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
