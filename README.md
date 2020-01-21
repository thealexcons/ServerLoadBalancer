# Server-side Load Balancer

A simple load balancer for servers/backends written in Go. 

## What is a load balancer?

Here is an in-depth [article](https://www.digitalocean.com/community/tutorials/what-is-load-balancing) of what it is, but essentially, it is a program that distributes the workload across different servers, for example, client requests. It is a key component in building scalable and performant services.

## Install

You can `go get github.com/thealexcons/server-load-balancer` to use this package in your own Go code.

## Usage
If you would like to use this package in your own code, take a look at `main.go` in the `example` directory for a self-explanatory and commented example usage of the package.

It can also be used as a command line tool, by building the `app.go` file in the `example` directory, by doing `go build app.go`.

I am in the process of building a CLI app for this load balancer, which I will link to when it is done. For now, you can use this by supplying two command line arguments: the path to a file containing the node addresses (and weights), and the port number for the load balancer to listen on:

`./app -nodes=nodes.txt -port=8080`

You can see the format of `nodes.txt` in the `example` directory.

## Some notes

The current code only serves HTTP requests, not HTTPS requests. So if you want to distribute requests to HTTPS endpoints, you must use the  `http.ListenAndServeTLS()` method, supplying an SSL certificate. I will probably add a feature that will allow an optional command line argument to an SSL certificate path some time soon, since it should be easy to do.

## Roadmap

- [ ] Add support for HTTPS traffic
- [ ] Add scheduling algorithm based on geolocation 
