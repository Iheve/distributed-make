package main

import (
	"fmt"
	. "github.com/Iheve/distributed-make/arith"
	"log"
	"net/rpc"
)

func main() {
	serverAddress := "localhost:1234"
	client, err := rpc.DialHTTP("tcp", serverAddress)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := Args{10, 4}
	// Asynchronous call
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	fmt.Println("result: ", replyCall.Reply)
}
