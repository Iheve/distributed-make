package main

import (
	. "github.com/Iheve/distributed-make/worker"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	worker := new(Worker)
	rpc.Register(worker)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
