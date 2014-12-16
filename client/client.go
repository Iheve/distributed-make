package main

import (
	"fmt"
	. "github.com/Iheve/distributed-make/worker"
	"log"
	"net/rpc"
)

func main() {
	serverAddress := "localhost:1234"
	client, err := rpc.DialHTTP("tcp", serverAddress)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := new(Args)
	args.Command = "date"
	// Asynchronous call
	response := new(Response)
	workerCall := client.Go("Worker.Output", args, response, nil)
	replyCall := <-workerCall.Done
	// check errors, print, etc.
	fmt.Println("result: ", replyCall.Reply)
}
