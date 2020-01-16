package main

import (
	"fmt"
	"net/http"
	"time"

	loadbalancer "github.com/thealexcons/ServerLoadBalancer"
)

func main() {
	// Create the server group and add nodes to it
	sg := &loadbalancer.ServerGroup{}
	sg.AddNode("localhost:4311", 1)
	// this node has weight 2 (will receive twice the no. of requests compared to other nodes)
	sg.AddNode("localhost:1232", 2)
	sg.AddNode("localhost:4192", 1)

	// Spin up the example nodes above
	go runExampleNodeServers()

	// Create the server at port 8080 and handle requests using the
	// load balanacer
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(sg.LoadBalancer),
	}

	// Start health checking routine every 2 minutes, with a timeout of 3 seconds
	sg.StartHealthChecker(time.Minute*2, time.Second*3)

	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Start some example nodes to simulate the load balancer
// Obviously, in practice, these nodes would be running separately.
func runExampleNodeServers() {
	server1 := &http.Server{
		Addr:    ":4311",
		Handler: http.HandlerFunc(exampleHandler),
	}
	server1.ListenAndServe()

	server2 := &http.Server{
		Addr:    ":1232",
		Handler: http.HandlerFunc(exampleHandler),
	}
	server2.ListenAndServe()

	server3 := &http.Server{
		Addr:    ":4192",
		Handler: http.HandlerFunc(exampleHandler),
	}
	server3.ListenAndServe()
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You reached node "+r.Host, nil)
}
