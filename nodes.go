package loadbalancer

import (
	"log"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Node struct {
	Address  *url.URL
	Weight   uint
	Alive    bool
	mutex    sync.RWMutex
	RevProxy *httputil.ReverseProxy
}

// SetAlive sets the node as alive or dead.
func (n *Node) SetAlive(alive bool) {
	n.mutex.Lock()
	n.Alive = alive
	n.mutex.Unlock()
}

// IsAlive returns true when node is alive
func (n *Node) IsAlive() bool {
	n.mutex.RLock()
	alive := n.Alive
	n.mutex.RUnlock()
	return alive
}

type ServerGroup struct {
	nodes   []*Node
	current uint32
	retries uint
}

func (s *ServerGroup) NextIndex() int {
	return int(atomic.AddUint32(&s.current, uint32(1))) % len(s.nodes)
}

// GetNextPeer returns next active peer to take a connection
// Algorithm here
func (s *ServerGroup) GetNextNode() *Node {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.nodes) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.nodes) // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if s.nodes[idx].IsAlive() {
			if i != next {
				atomic.StoreUint32(&s.current, uint32(idx)) // mark the current one
			}
			return s.nodes[idx]
		}
	}
	return nil
}

// AddNode creates a node and adds it to the server group given a nodeUrl and weight
func (s *ServerGroup) AddNode(nodeUrl string, weight uint) {
	addr, err := url.Parse(nodeUrl)
	if err != nil {
		log.Fatal("The URL provided '%s' is not a valid address.", nodeUrl)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(addr)
	proxy.ErrorHandler = s.proxyErrorHandler
	node := &Node{
		Address:  addr,
		Weight:   weight,
		Alive:    true,
		RevProxy: proxy,
	}
	s.nodes = append(s.nodes, node)
}

// CheckAndUpdateHealth sets and reports back which backends are alive
func (s *ServerGroup) CheckAndUpdateHealth() {
	for _, n := range s.nodes {
		alive := is
	}
}
