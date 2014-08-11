gogo-proxy
==========

A http &amp; websocket reverse proxy written in Go


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
