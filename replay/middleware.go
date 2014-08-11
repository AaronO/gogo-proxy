package replay

import (
	"net/http"
	"time"
)

// Middleware allows us to wrap http.Handler to be retried
// this is useful for building robust proxies
type Middleware struct {
	retries      int
	period       time.Duration
	handler      http.Handler
	ErrorHandler func(rw http.ResponseWriter, req *http.Request, err error)
}

// NewMiddleware
func NewMiddleware(retries int, period time.Duration, handler http.Handler) *Middleware {
	return &Middleware{
		retries: retries,
		period:  period,
		handler: handler,
	}
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	replayer, err := NewReplayer(req, rw)
	if err != nil {
		return
	}

	err = retryWait(func() error {
		newReq, err := replayer.Replay()
		if err != nil {
			return err
		}

		m.handler.ServeHTTP(replayer, newReq)

		if err := replayer.GetError(); err != nil {
			replayer.Stop()
			return err
		}

		return nil
	}, m.retries, m.period)

	if err != nil {
		if m.ErrorHandler != nil {
			m.ErrorHandler(rw, req, err)
		} else {
			defaultErrorHandler(rw, req, err)
		}
	}
}

func defaultErrorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}
