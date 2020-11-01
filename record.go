package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

// Record stores the request's timestamp corresponding to a certain IP.
type Record struct {
	mtx       sync.Mutex
	timestamp []int64
}

var LimitExceedError error = errors.New("Rate limit exceeded")

// NewRecord returns a pointer of Record struct.
func NewRecord() *Record {
	return &Record{
		timestamp: make([]int64, 0),
	}
}

// Add adds a record if it is allowed regards to the condition (limit and window),
// return LimitExceedError if the number of requests associating to the incoming IP
// exceeds the limit.
//
// Add serves each request synchronously to avoid race condition.
func (rec *Record) Add(limit, window int) (int, error) {
	rec.mtx.Lock()
	defer rec.mtx.Unlock()

	now := time.Now().Unix()
	start := now - int64(window)

	l := 0
	r := len(rec.timestamp) - 1
	for l <= r {
		m := l + (r-l)/2
		if rec.timestamp[m] > start {
			r = m - 1
		} else {
			l = m + 1
		}
	}

	// l will be the first index that is greater than start
	cnt := len(rec.timestamp) - l

	// lazy cleanup old reacord and add new record.
	rec.timestamp = rec.timestamp[l:]

	if cnt >= limit {
		return cnt, LimitExceedError
	}

	// it is allowed to add record.
	rec.timestamp = append(rec.timestamp, now)

	return cnt + 1, nil
}
