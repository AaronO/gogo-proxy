package replay

import (
	"fmt"
	"net/http"
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

func NewReplayer(req *http.Request, w http.ResponseWriter) (*Replayer, error) {
	r := Replayer{
		req:    req,
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
	return r.writer.Header()
}

func (r *Replayer) Write(data []byte) (int, error) {
	// Don't write any data (this will be retried)
	if r.IsFailed() {
		return len(data), nil
	}

	r.Play.Bytes += len(data)
	r.Play.Writes += 1
	return r.writer.Write(data)
}

func (r *Replayer) WriteHeader(status int) {
	r.Play.Status = status

	// This request is failing don't write out please ...
	if r.IsFailed() {
		return
	}

	r.writer.WriteHeader(status)
}

func (r *Replayer) GetError() error {
	if r.Play.Status >= 500 {
		return fmt.Errorf("Status code is %d", r.Play.Status)
	}

	return nil
}

func (r *Replayer) IsFailed() bool {
	return r.GetError() != nil
}
