package main

import (
	"fmt"
	. "github.com/Iheve/distributed-make/arith"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
	fmt.Printf("Hello, world.\n")
}
