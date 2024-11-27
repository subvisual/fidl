package proxy

import (
	"context"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type ctxKey int

type UpstreamTracker struct {
	mu       sync.Mutex
	trackers map[string]*int64
}

type TrackingReader struct {
	io.ReadCloser
	tracker   *UpstreamTracker
	trackerID string
}

const trackerCtxKey ctxKey = iota

func NewUpstreamTracker() *UpstreamTracker {
	return &UpstreamTracker{
		trackers: make(map[string]*int64),
	}
}

func (t *UpstreamTracker) Start(r *http.Request) (*int64, *http.Request, func()) {
	t.mu.Lock()
	defer t.mu.Unlock()

	trackerID := uuid.New().String()
	accumulator := new(int64)
	t.trackers[trackerID] = accumulator
	ctx := context.WithValue(r.Context(), trackerCtxKey, trackerID)
	cleanupFn := func() { t.cleanup(trackerID) }

	return accumulator, r.WithContext(ctx), cleanupFn
}

func (t *UpstreamTracker) Track(trackerID string, n int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if counter, ok := t.trackers[trackerID]; ok {
		atomic.AddInt64(counter, int64(n))
	}
}

func (t *UpstreamTracker) cleanup(trackerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.trackers, trackerID)
}

func (r *TrackingReader) Read(p []byte) (int, error) {
	bytesRead, err := r.ReadCloser.Read(p)
	if bytesRead > 0 {
		r.tracker.Track(r.trackerID, bytesRead)
	}

	// nolint:wrapcheck
	return bytesRead, err
}
