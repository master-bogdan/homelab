// Package metrics provides a minimal Prometheus-compatible metrics registry.
package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var defaultRegistry = newRegistry()

var httpDurationBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

type registry struct {
	mu                  sync.RWMutex
	httpRequests        map[httpRequestKey]uint64
	httpDuration        map[httpDurationKey]*histogram
	roomLifecycleCounts map[string]uint64
	wsActiveConnections int64
	wsTotalConnections  uint64
}

type httpRequestKey struct {
	Method string
	Path   string
	Status int
}

type httpDurationKey struct {
	Method string
	Path   string
}

type histogram struct {
	Count   uint64
	Sum     float64
	Buckets []uint64
}

func newRegistry() *registry {
	return &registry{
		httpRequests:        make(map[httpRequestKey]uint64),
		httpDuration:        make(map[httpDurationKey]*histogram),
		roomLifecycleCounts: make(map[string]uint64),
	}
}

func RecordHTTPRequest(method, path string, status int, duration time.Duration) {
	defaultRegistry.recordHTTPRequest(method, path, status, duration)
}

func IncWSConnections() {
	defaultRegistry.incWSConnections()
}

func DecWSConnections() {
	defaultRegistry.decWSConnections()
}

func RecordRoomLifecycle(state string) {
	defaultRegistry.recordRoomLifecycle(state)
}

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		_, _ = w.Write([]byte(defaultRegistry.render()))
	})
}

func (r *registry) recordHTTPRequest(method, path string, status int, duration time.Duration) {
	method = normalizeLabelValue(method)
	path = normalizePath(path)
	if status <= 0 {
		status = http.StatusOK
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.httpRequests[httpRequestKey{
		Method: method,
		Path:   path,
		Status: status,
	}]++

	histKey := httpDurationKey{
		Method: method,
		Path:   path,
	}
	hist := r.httpDuration[histKey]
	if hist == nil {
		hist = &histogram{Buckets: make([]uint64, len(httpDurationBuckets))}
		r.httpDuration[histKey] = hist
	}

	seconds := duration.Seconds()
	hist.Count++
	hist.Sum += seconds
	for i, bucket := range httpDurationBuckets {
		if seconds <= bucket {
			hist.Buckets[i]++
			break
		}
	}
}

func (r *registry) incWSConnections() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.wsActiveConnections++
	r.wsTotalConnections++
}

func (r *registry) decWSConnections() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.wsActiveConnections > 0 {
		r.wsActiveConnections--
	}
}

func (r *registry) recordRoomLifecycle(state string) {
	normalized := normalizeLabelValue(state)
	if normalized == "" {
		return
	}

	r.mu.Lock()
	r.roomLifecycleCounts[normalized]++
	r.mu.Unlock()
}

func (r *registry) render() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var b strings.Builder

	b.WriteString("# HELP http_requests_total Total number of HTTP requests.\n")
	b.WriteString("# TYPE http_requests_total counter\n")
	requestKeys := make([]httpRequestKey, 0, len(r.httpRequests))
	for key := range r.httpRequests {
		requestKeys = append(requestKeys, key)
	}
	sort.Slice(requestKeys, func(i, j int) bool {
		if requestKeys[i].Method != requestKeys[j].Method {
			return requestKeys[i].Method < requestKeys[j].Method
		}
		if requestKeys[i].Path != requestKeys[j].Path {
			return requestKeys[i].Path < requestKeys[j].Path
		}
		return requestKeys[i].Status < requestKeys[j].Status
	})
	for _, key := range requestKeys {
		fmt.Fprintf(
			&b,
			"http_requests_total{method=%q,path=%q,status=%q} %d\n",
			key.Method,
			key.Path,
			strconv.Itoa(key.Status),
			r.httpRequests[key],
		)
	}

	b.WriteString("# HELP http_request_duration_seconds HTTP request duration in seconds.\n")
	b.WriteString("# TYPE http_request_duration_seconds histogram\n")
	durationKeys := make([]httpDurationKey, 0, len(r.httpDuration))
	for key := range r.httpDuration {
		durationKeys = append(durationKeys, key)
	}
	sort.Slice(durationKeys, func(i, j int) bool {
		if durationKeys[i].Method != durationKeys[j].Method {
			return durationKeys[i].Method < durationKeys[j].Method
		}
		return durationKeys[i].Path < durationKeys[j].Path
	})
	for _, key := range durationKeys {
		hist := r.httpDuration[key]
		var cumulative uint64
		for i, bucket := range httpDurationBuckets {
			cumulative += hist.Buckets[i]
			fmt.Fprintf(
				&b,
				"http_request_duration_seconds_bucket{method=%q,path=%q,le=%q} %d\n",
				key.Method,
				key.Path,
				strconv.FormatFloat(bucket, 'f', -1, 64),
				cumulative,
			)
		}
		fmt.Fprintf(
			&b,
			"http_request_duration_seconds_bucket{method=%q,path=%q,le=%q} %d\n",
			key.Method,
			key.Path,
			"+Inf",
			hist.Count,
		)
		fmt.Fprintf(
			&b,
			"http_request_duration_seconds_sum{method=%q,path=%q} %s\n",
			key.Method,
			key.Path,
			strconv.FormatFloat(hist.Sum, 'f', 6, 64),
		)
		fmt.Fprintf(
			&b,
			"http_request_duration_seconds_count{method=%q,path=%q} %d\n",
			key.Method,
			key.Path,
			hist.Count,
		)
	}

	b.WriteString("# HELP ws_connections_active Current number of active websocket connections.\n")
	b.WriteString("# TYPE ws_connections_active gauge\n")
	fmt.Fprintf(&b, "ws_connections_active %d\n", r.wsActiveConnections)

	b.WriteString("# HELP ws_connections_total Total number of websocket connections accepted.\n")
	b.WriteString("# TYPE ws_connections_total counter\n")
	fmt.Fprintf(&b, "ws_connections_total %d\n", r.wsTotalConnections)

	b.WriteString("# HELP room_lifecycle_total Total number of room lifecycle transitions.\n")
	b.WriteString("# TYPE room_lifecycle_total counter\n")
	roomStates := make([]string, 0, len(r.roomLifecycleCounts))
	for state := range r.roomLifecycleCounts {
		roomStates = append(roomStates, state)
	}
	sort.Strings(roomStates)
	for _, state := range roomStates {
		fmt.Fprintf(&b, "room_lifecycle_total{state=%q} %d\n", state, r.roomLifecycleCounts[state])
	}

	return b.String()
}

func normalizePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/"
	}
	if !strings.HasPrefix(trimmed, "/") {
		return "/" + trimmed
	}
	return trimmed
}

func normalizeLabelValue(value string) string {
	return strings.TrimSpace(strings.ToUpper(value))
}
