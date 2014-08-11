package proxy

import (
	"reflect"
	"testing"

	"github.com/AaronO/gogo-proxy"
)

func TestRoundrobinEmpty(t *testing.T) {
	balancer := proxy.RoundrobinBalancer()

	// Getting a backend from an empty round robin balancer should fail
	if _, err := balancer(nil); err == nil {
		t.Fatalf("RoundrobinBalancer should fail when given no hosts")
	}
}

func TestRandomEmpty(t *testing.T) {
	balancer := proxy.RandomBalancer()

	// Getting a backend from an empty random balancer should fail
	if _, err := balancer(nil); err == nil {
		t.Fatalf("RandomBalancer should fail when given no hosts")
	}
}

func TestRoundrobin(t *testing.T) {
	balancer := proxy.RoundrobinBalancer("1", "2", "3")

	hosts := []string{}

	// loop twice + once more
	for i := 0; i < 7; i++ {
		host, err := balancer(nil)
		if err != nil {
			t.Error(err)
		}
		hosts = append(hosts, host)
	}

	equal := reflect.DeepEqual(
		hosts,
		[]string{"1", "2", "3", "1", "2", "3", "1"},
	)

	if !equal {
		t.Fatalf("RoundrobinBalancer did not generate expected hosts")
	}
}
