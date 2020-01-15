package main

import "github.com/thealexcons/ServerLoadBalancer"

func main() {
	// Create the server group and add nodes to it
	sg := &ServerGroup{}
	sg.AddNode("localhost:4311", 1)
	// this node has weight 2 (will receive twice the no. of requests compared to other nodes)
	sg.AddNode("localhost:1232", 2) 
	sg.AddNode("localhost:4192", 1)

	// Create the server at port 8080 and handle requests using the 
	// load balanacer
	server := &http.Server{
		Addr: 	 "8080",
		Handler: http.HandlerFunc(sg.LoadBalancer)
	}

	// Start health checking routine every 2 minutes, with a timeout of 3 seconds
	sg.StartHealthChecker(time.Minute * 2, time.Second * 3)

	// Start the server
	_ := server.ListenAndServe()
}
