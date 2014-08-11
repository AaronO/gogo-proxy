package proxy

import (
	"errors"
	"math/rand"
	"net/http"
)

// RoundrobinBalancer generates a balancer that
// cycles through the list of hosts in a round robin fashion (equally distributing traffic)
func RoundrobinBalancer(hosts ...string) func(*http.Request) (string, error) {
	n := uint32(len(hosts))
	var idx uint32 = 0
	return func(req *http.Request) (string, error) {
		if n == 0 {
			return "", errors.New("RoundrobinBalancer can not work on an empty list")
		}

		h := hosts[idx%n]

		// Increment idx to add rotation
		idx++

		return h, nil
	}
}

// RandomBalancer generates a balancer than
// randomly picks one of provided hosts
// (this does not guarantee equal traffic distribution) but with many requests should get close to it
func RandomBalancer(hosts ...string) func(*http.Request) (string, error) {
	n := len(hosts)
	return func(req *http.Request) (string, error) {
		if n == 0 {
			return "", errors.New("RandomBalancer can not work on an empty list")
		}

		idx := rand.Intn(n)
		return hosts[idx], nil
	}
}
