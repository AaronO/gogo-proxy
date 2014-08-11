gogo-proxy
==========

A fast and robust http/websocket reverse proxy written in Go


### Features
  - **Simple:**
    - Implements `http.Handler` interface
    - Easy to write `Balancer`s for custom routing logic
    - Built-in balancing patterns (`Roundrobin`, `Random`, etc ...)
  - **Robust:** Retry requests on failure
  - **Flexible:**
    - Custom error handling (so you can draw custom error pages etc ...) (use `ErrorHandler`)
    - Custom request rewriting (use `Rewriter`)
  - **Fast & "scalable":**
    - Written in go, so concurrent by default and fast


### Example

```go
package main

import (
    "github.com/AaronO/gogo-proxy"
    "net/http"
)

func main() {
    p, _ := proxy.New(proxy.ProxyOptions{
        Balancer: func(req *http.Request) (string, error) {
            return "https://www.google.com", nil
        },
    })

    http.ListenAndServe(":8080", p)
}
```
