package loadbalancer

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	Attempts int = iota
	Retry
)

// LoadBalancer starts load balancing
func (s *ServerGroup) LoadBalancer(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "The service is unavailable.", http.StatusServiceUnavailable)
	}

	node := s.GetNextNode()
	if node != nil {
		node.RevProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "The service is unavailable.", http.StatusServiceUnavailable)
}

func CheckHealth(refreshRate time.Duration) {
	clock := time.NewTicker(refreshRate)
	for {
		select {
		case <-clock.C:

		}
	}
}

// GetAttemptsFromContext returns the number of attempts for a given request
func GetAttemptsFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// Error handler for the proxy.
func (s *ServerGroup) proxyErrorHandler(w http.ResponseWriter, r *http.Request, e error) {
	log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
	retries := GetRetryFromContext(r)
	if retries < 3 {
		select {
		case <-time.After(10 * time.Millisecond):
			ctx := context.WithValue(request.Context(), Retry, retries+1)
			proxy.ServeHTTP(w, r.WithContext(ctx))
		}
		return
	}

	// after 3 retries, mark this backend as down
	s.MarkBackendStatus(serverUrl, false)

	// if the same request routing for few attempts with different backends, increase the count
	attempts := GetAttemptsFromContext(r)
	log.Printf("%s(%s) Attempting retry %d\n", r.RemoteAddr, r.URL.Path, attempts)
	ctx := context.WithValue(r.Context(), Attempts, attempts+1)
	lb(w, r.WithContext(ctx))
}
