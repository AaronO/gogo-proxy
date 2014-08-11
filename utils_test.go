package proxy

import (
	"testing"
)

func TestUrlValidation(t *testing.T) {
	if err := validateUrl(""); err == nil {
		t.Errorf("Empty URLs should not pass validation")
	}

	if err := validateUrl("http://"); err == nil {
		t.Errorf("URLs with only a protocol should not pass validation")
	}

	if err := validateUrl("http://127.0.0.1:5000/"); err != nil {
		t.Errorf("Full URLs should pass validation")
	}
}

func TestUrlNormalization(t *testing.T) {
	shouldEqual(normalizeUrl(""), "", t)
	shouldEqual(normalizeUrl("://"), "", t)
	shouldEqual(normalizeUrl("http://"), "", t)
	shouldEqual(normalizeUrl("https://"), "", t)

	shouldEqual(normalizeUrl("192.168.1.1"), "http://192.168.1.1/", t)
	shouldEqual(normalizeUrl("://192.168.1.1"), "http://192.168.1.1/", t)
	shouldEqual(normalizeUrl("http://192.168.1.1"), "http://192.168.1.1/", t)
	shouldEqual(normalizeUrl("https://192.168.1.1"), "https://192.168.1.1/", t)

	shouldEqual(normalizeUrl("example.com"), "http://example.com/", t)
	shouldEqual(normalizeUrl("example.com"), "http://example.com/", t)
	shouldEqual(normalizeUrl("example.com/dir/"), "http://example.com/dir/", t)

	shouldEqual(normalizeUrl("localhost:3000"), "http://localhost:3000/", t)
	shouldEqual(normalizeUrl("example.com:3000"), "http://example.com:3000/", t)
	shouldEqual(normalizeUrl("https://example.com:3000"), "https://example.com:3000/", t)
	shouldEqual(normalizeUrl("192.168.1.1:3000"), "http://192.168.1.1:3000/", t)
}

func TestHttpScheme(t *testing.T) {
	shouldEqual(httpScheme("http"), "http", t)
	shouldEqual(httpScheme("https"), "https", t)
	shouldEqual(httpScheme("ws"), "http", t)
	shouldEqual(httpScheme("wss"), "https", t)

	shouldEqual(httpScheme(""), "http", t)
	shouldEqual(httpScheme("abcd"), "http", t)
}

func TestWebsocketScheme(t *testing.T) {
	shouldEqual(websocketScheme("ws"), "ws", t)
	shouldEqual(websocketScheme("wss"), "wss", t)
	shouldEqual(websocketScheme("https"), "wss", t)
	shouldEqual(websocketScheme("https"), "wss", t)

	shouldEqual(websocketScheme(""), "ws", t)
	shouldEqual(websocketScheme("abcd"), "ws", t)
}

// Utility function for shorter string compares
func shouldEqual(value, expected string, t *testing.T) {
	if value != expected {
		t.Errorf("Value is '%s', expected '%s'", value, expected)
	}
}
