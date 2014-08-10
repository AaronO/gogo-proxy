package proxy_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/AaronO/gogo-proxy"
)

func TestGoogle(t *testing.T) {
	// Create proxy to google.co.uk
	p, _ := proxy.New(proxy.ProxyOptions{
		Balancer: func(req *http.Request) (string, error) {
			return "https://www.google.co.uk" + req.URL.Path, nil
		},
	})

	ch := make(chan error, 2)

	go func() {
		ch <- http.ListenAndServe(":57439", p)
	}()

	res, err := http.Get("http://localhost:57439/")
	if err != nil {
		t.Fatal(err)
	}

	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(page, []byte("I'm Feeling Lucky")) {
		t.Fatalf("Proxied Google homepage does not have lucky button !")
	}
}
