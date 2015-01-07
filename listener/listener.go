package main

import (
	. "github.com/Iheve/distributed-make/worker"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func main() {
	var port string

	worker := new(Worker)
	rpc.Register(worker)
	rpc.HandleHTTP()

	// Check the args, use default port (= 4242)
	if len(os.Args) != 2 {
		port = ":4242"
	} else {
		port = ":" + os.Args[1]
	}

	l, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
