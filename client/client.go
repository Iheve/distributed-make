package main

import (
	"fmt"
	"github.com/Iheve/distributed-make/parser"
	"github.com/Iheve/distributed-make/worker"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

func run(client *rpc.Client, todo chan *parser.Task) {
	for {
		t := <-todo
		var response worker.Response
		args := new(worker.Args)
		fmt.Println("Target :", t.Target)
		args.Target = t.Target
		args.Cmds = t.Cmds
		//Pack dependencies
		for _, d := range t.Deps {
			fmt.Printf("Target: %s Dep: %s\n", t.Target, d)
			if d == "" {
				continue
			}
			var f worker.File
			f.Name = d
			info, _ := os.Stat(d)
			f.Mode = info.Mode()
			var err error
			f.Content, err = ioutil.ReadFile(d)
			if err != nil {
				log.Fatal("Cant read file: ", d, " : ", err)
			}
			args.Deps = append(args.Deps, f)
		}
		//Synchronous call TODO switch to asynchronous ?
		err := client.Call("Worker.Output", args, &response)
		if err != nil {
			log.Fatal("RPC call error:", err)
		}
		//Unpack target
		err = ioutil.WriteFile(response.Target.Name, response.Target.Content, response.Target.Mode)
		if err != nil {
			log.Fatal("Can not create file: ", response.Target.Name, " : ", err)
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
