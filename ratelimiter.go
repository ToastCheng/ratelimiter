package ratelimiter

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// RateLimitHandler return the original http.Handler with a midleware.
func RateLimitHandler(h http.Handler, limit, window int) http.Handler {
	r := rateLimitHandler{
		records: map[string]*Record{},
		limit:   limit,
		window:  window,
		handler: h,
	}
	return r
}

// rateLimitHandler implement ServeHTTP,
// handles the incoming request and performs rate limiting.
type rateLimitHandler struct {
	records map[string]*Record
	limit   int
	window  int
	handler http.Handler
}

// ServeHTTP serves the http request.
func (h rateLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// get the IP address of the request.
	ip, err := h.getIP(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to decode ip: %v", err)
		return
	}
	log.Printf("get request from: %s", ip)

	// if the ip address has no record yet, initialize a Record for it.
	if _, exists := h.records[ip]; !exists {
		h.records[ip] = NewRecord()
	}

	// try to add the record.
	_, err = h.records[ip].Add(h.limit, h.window)
	if err == LimitExceedError {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, "Error")
		return
	}

	h.handler.ServeHTTP(w, r)
}

func (h rateLimitHandler) getIP(r *http.Request) (string, error) {
	// if the request came from a reverse proxy such as Nginx,
	// use the information in X-Forwarded-For header.
	if r.Header.Get("X-Forwarded-For") != "" {
		// TODO: Assume X-Forwarded-For contains single IP for now, but it might be a list.
		return r.Header.Get("X-Forwarded-For"), nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}
