package main

import (
	"flag"
	"github.com/Iheve/distributed-make/worker"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	var help bool
	var port string
	flag.BoolVar(&help, "help", false, "Display this helper message")
	flag.StringVar(&port, "port", "4242", "Port of the listener")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	worker := new(worker.Worker)
	rpc.Register(worker)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Println("Listening on port " + port)
	http.Serve(l, nil)
}
