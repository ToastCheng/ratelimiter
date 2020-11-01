# Rate Limiter

A middleware to rate-limit the number of requests down to `limit` within `window` second for each IP address.
If a certain IP requests exceeded the `limit` within `window` second, it will response with `Too Many Request (429)`.

### Usage

```go
import (
    "net/http"
    "github.com/toastcheng/ratelimiter"
)

func main() {
    r := http.NewServeMux()

    limit := 60
    window := 60
    
    r.Handle("/web", ratelimiter.RateLimitHandler(handler, limit, window))
    http.ListenAndServe(":8000", r)
}

type handler struct {}
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // ...
}
```