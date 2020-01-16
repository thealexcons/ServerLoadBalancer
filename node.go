package loadbalancer

import (
	"log"
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type node struct {
	address  *url.URL
	weight   int
	alive    bool
	mutex    sync.RWMutex
	revProxy *httputil.ReverseProxy
}

// SetAlive sets the node as alive or dead.
func (n *node) setAlive(alive bool) {
	n.mutex.Lock()
	n.alive = alive
	n.mutex.Unlock()
}

// IsAlive returns true when node is alive
func (n *node) isAlive() bool {
	n.mutex.RLock()
	alive := n.alive
	n.mutex.RUnlock()
	return alive
}

// isNodeAlive checks whether a node is alive by establishing a TCP connection
func (n *node) checkIfAlive(timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", n.address.Host, timeout)
	defer conn.Close()
	if err != nil {
		log.Println("Node unreachable, error: ", err)
		return false
	}
	return true
}
