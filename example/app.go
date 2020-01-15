package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thealexcons/ServerLoadBalancer"
)

func main() {
	var nodeDataPath string
	var port int
	flag.StringVar(&nodeDataPath, "nodes", "", "Load nodes file path, containing weights if required")
	flag.IntVar(&port, "port", 8080, "Port to serve. Default is 8080")
	flag.Parse()

	// Read data file
	data, err := ioutil.ReadFile(nodeDataPath)
	if err != nil {
		log.Fatal("Error reading the provided file, " + err.Error())
	}
	nodes := strings.Split(string(data), "\n")

	sg := &ServerGroup{}
	sg.SetRetries(3)

	// Add nodes from file to ServerGroup
	for _, nd := range nodes {
		nodeData := strings.Split(nd, ",")
		nodeAddr := nodeData[0]
		if len(nodeData) == 1 {
			sg.AddNode(nodeAddr, 1) // no weight provided, set to 1
			continue
		} else if len(nodeData) == 2 {
			nodeWeightStr := strings.Split(nd, ",")[1]
			nodeWeight, err := strconv.Atoi(nodeWeightStr)
			if err != nil {
				log.Fatal("The port %s is not valid.", nodeWeightStr)
			}
			sg.AddNode(nodeAddr, nodeWeight)
		} else {
			log.Fatal("Please provided a valid data file")
		}

	}

	server := &http.Server{
		Addr:    string(port),
		Handler: http.HandlerFunc(sg.LoadBalancer),
	}

	// Start health checking routine every 2 minutes, with a timeout of 3 seconds
	sg.StartHealthChecker(time.Minute*2, time.Second*3)

	// Start the server
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Could not start server: ", err.Error())
	}

}
