# Rate Limiter

### Usage

```go
import (
    "net/http"
    "github.com/toastcheng/ratelimiter"
)

func main() {
    r := http.NewServeMux()

    r.Handle("/web", ratelimiter.RateLimitHandler(handler, 60, 60))

}

type handler struct {}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // ...
}
```