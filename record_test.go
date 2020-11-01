package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecord(t *testing.T) {
	r := NewRecord()
	assert.Equal(t, 0, len(r.timestamp))
}

func TestAdd(t *testing.T) {
	r := NewRecord()
	for i := 0; i < 60; i++ {
		cnt, err := r.Add(60, 60)
		assert.NoError(t, err)
		assert.Equal(t, i+1, cnt)
	}
	_, err := r.Add(60, 60)
	if assert.Error(t, err) {
		assert.Equal(t, LimitExceedError, err)
	}

}
