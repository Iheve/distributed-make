package main

import (
	"flag"
	"fmt"
	"github.com/Iheve/distributed-make/config"
	"github.com/Iheve/distributed-make/parser"
	"github.com/Iheve/distributed-make/worker"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

func run(client *rpc.Client, name string, todo chan *parser.Task) {
	for {
		t := <-todo
		var response worker.Response
		args := new(worker.Args)
		log.Println(name, " target:", t.Target)
		args.Target = t.Target
		args.Cmds = t.Cmds
		//Pack dependencies
		for _, d := range t.Deps {
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
        /*
		fmt.Println("Command done, outputs:")
		for _, s := range response.Output {
			fmt.Print(s)
		}
        */
	}
}

func walk(t *parser.Task, todo chan *parser.Task, depth int) bool {
    /*
    for i:=0; i < depth; i++ {
        fmt.Print("\t")
    }
    fmt.Print(t.Target, ":", t.Done, "\n")
    */
	if t.Done {
		return true
	}

	if t.Affected {
		return false
	}

	res := true
	for _, s := range t.Sons {
		if s != nil {
			res = walk(s, todo, depth + 1) && res
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

	hostFilePtr := flag.String("hostfile", "hostfile.cfg", "a string")
	flag.Parse()

	fmt.Println("Hostfile: ", *hostFilePtr)
	fmt.Println("Args: ", flag.Args())

	args := flag.Args()

	if len(args) != 1 {
		path = "Makefile"
	} else {
		path = args[0]
	}

	head, err := parser.Parse(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	parser.Print(head, 0)
	hosts := config.Parse(*hostFilePtr)

	fmt.Println("Hosts:", hosts)
	todo := make(chan *parser.Task) //TODO set the buffer lenght in function of the number of worker

	for i := range hosts {
		serverAddress := hosts[i]
		client, err := rpc.DialHTTP("tcp", serverAddress)
		if err != nil {
			log.Fatal("dialing:", err)
		}

		go run(client, serverAddress, todo) //TODO run for each worker
	}

	for !walk(head, todo, 0) {
	}

}
