package main

import (
	"fmt"
	"github.com/Iheve/distributed-make/parser"
	"github.com/Iheve/distributed-make/worker"
	"log"
	"net/rpc"
	"os"
)

func run(client *rpc.Client, todo chan *parser.Task) {
	for {
		t := <-todo
		var response worker.Response
		//Synchronous call TODO switch to asynchronous ?
		args := new(worker.Args)
		args.Cmds = t.Cmds
		err := client.Call("Worker.Output", args, &response)
		if err != nil {
			log.Fatal("RPC call error:", err)
		}
		t.Done = true
		fmt.Println("Command done, outputs:")
		for _, s := range response.Output {
			fmt.Print(s)
		}
	}
}

func walk(t *parser.Task, todo chan *parser.Task) bool {
	if t.Done {
		return true
	}

	if t.Affected {
		return false
	}

	res := true
	for _, s := range t.Sons {
		if s != nil {
			res = res && walk(s, todo)
		}
	}

	if res {
		t.Affected = true
		todo <- t
	}

	return false
}

func main() {

	var path string
	if len(os.Args) != 2 {
		path = "Makefile"
	} else {
		path = os.Args[1]
	}

	head, err := parser.Parse(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	parser.Print(head, 0)

	serverAddress := "localhost:1234"
	client, err := rpc.DialHTTP("tcp", serverAddress)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	todo := make(chan *parser.Task, 10) //TODO set the buffer lenght in function of the number of worker

	go run(client, todo) //TODO run for each worker
	for !walk(head, todo) {
	}

}
