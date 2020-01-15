package loadbalancer

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type ServerGroup struct {
	nodes   []*node
	current uint32
	retries int // default is 3 retries
}

func (s *ServerGroup) nextIndex() int {
	return int(atomic.AddUint32(&s.current, uint32(1))) % len(s.nodes)
}

// GetNextPeer returns next active peer to take a connection
// Uses Weighted Round Robin algorithm
func (s *ServerGroup) getNextNode() *node {
	/*
		// loop entire backends to find out an Alive backend
		next := s.nextIndex()
		l := len(s.nodes) + next // start from next and move a full cycle
		for i := next; i < l; i++ {
			idx := i % len(s.nodes) // take an index by modding with length
			// if we have an alive backend, use it and store if its not the original one
			if s.nodes[idx].isAlive() {
				if i != next {
					atomic.StoreUint32(&s.current, uint32(idx)) // mark the current one
				}
				return s.nodes[idx]
			}
		}
		return nil
	*/
}

// markBackendStatus changes a alive status of a node
func (s *ServerGroup) markBackendStatus(nodeUrl *url.URL, alive bool) {
	for _, n := range s.nodes {
		if n.address.String() == nodeUrl.String() {
			n.setAlive(alive)
			break
		}
	}
}

// AddNode creates a node and adds it to the server group given a nodeUrl and weight
func (s *ServerGroup) AddNode(nodeUrl string, weight uint) {
	addr, err := url.Parse(nodeUrl)
	if err != nil {
		log.Fatal("The URL provided '%s' is not a valid address.", nodeUrl)
		return
	}

	// Create reverse proxy and define its ErrorHandler
	proxy := httputil.NewSingleHostReverseProxy(addr)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("[%s] %s\n", addr.Host, e.Error())
		retries := GetRetryFromContext(r)
		if retries < s.retries {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(r.Context(), Retry, retries+1)
				proxy.ServeHTTP(w, r.WithContext(ctx))
			}
			return
		}

		// after n retries, mark this backend as down
		s.markBackendStatus(addr, false)

		attempts := GetAttemptsFromContext(r)
		log.Printf("%s(%s) Retrying %d\n", r.RemoteAddr, r.URL.Path, attempts)
		ctx := context.WithValue(r.Context(), Attempts, attempts+1)
		s.LoadBalancer(w, r.WithContext(ctx))
	}

	node := &node{
		address:  addr,
		weight:   weight,
		alive:    true,
		revProxy: proxy,
	}

	s.nodes = append(s.nodes, node)
}

// LoadBalancer starts load balancing the server group
func (s *ServerGroup) LoadBalancer(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > s.retries {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "The service is unavailable.", http.StatusServiceUnavailable)
	}

	node := s.getNextNode()
	if node != nil {
		node.revProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "The service is unavailable.", http.StatusServiceUnavailable)
}

// CheckAndUpdateHealth sets and reports back which backends are alive
func (s *ServerGroup) checkAndUpdateHealth(timeout time.Duration) {
	for _, n := range s.nodes {
		status := "alive"
		alive := n.checkIfAlive(timeout)
		n.setAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("Node %s is %s\n", n.address, status)
	}
}

// StartHealthChecker launches a go routine that checks the health.
// refreshRate is how often each node will be checked and timeout is the
// max. time to wait for a response.
func (s *ServerGroup) StartHealthChecker(refreshRate time.Duration, timeout time.Duration) {
	if timeout >= refreshRate {
		log.Fatal("Timeout should be less than refresh rate when checking node health")
	}
	go func() {
		clock := time.NewTicker(refreshRate)
		for {
			select {
			case <-clock.C:
				s.checkAndUpdateHealth(timeout)
			}
		}
	}()
}
