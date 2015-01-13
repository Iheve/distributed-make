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

func run(host string, todo chan *parser.Task, verbose bool) {
	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		log.Println("Can not contact", host, err)
		return
	}
	for {
		t := <-todo
		var response worker.Response
		args := new(worker.Args)
		log.Println(host, "builds", t.Target)
		args.Target = t.Target
		args.Cmds = t.Cmds
		//Pack dependencies
		for _, d := range t.Deps {
			if d == "" {
				continue
			}
			var f worker.File
			f.Name = d
			var err error
			f.Content, err = ioutil.ReadFile(d)
			if err != nil {
				log.Fatal("Cant read file: ", d, " : ", err)
			}
			info, _ := os.Stat(d)
			f.Mode = info.Mode()
			args.Deps = append(args.Deps, f)
		}
		//Synchronous call
		err := client.Call("Worker.Output", args, &response)
		if err != nil {
			s := fmt.Sprintf("%v", err)
			if s == "unexpected EOF" || s == "connection is shut down" {
				log.Println("Contact lost with ", host)
				log.Println(t.Target, "will be rebuilt.")
				log.Println(host, "will not receive job anymore.")
				todo <- t
				return
			}
			log.Fatal(host, " failed to build target ", t.Target, ":", err)
		}
		//Unpack target
		err = ioutil.WriteFile(response.Target.Name, response.Target.Content, response.Target.Mode)
		if err != nil {
			log.Fatal("Can not create file: ", response.Target.Name, " : ", err)
		}

		t.Done = true
		if verbose {
			log.Println(host, "Command done, outputs:")
			for _, s := range response.Output {
				log.Println("\n", s)
			}
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
			res = walk(s, todo) && res
		}
	}

	if res {
		t.Affected = true
		todo <- t
	}

	return false
}

func main() {
	var help, verbose, showGraph bool
	var hostfileName, makefileName string
	flag.BoolVar(&help, "help", false, "Display this helper message")
	flag.BoolVar(&verbose, "verbose", false, "Show outputs of commands")
	flag.BoolVar(&showGraph, "showgraph", false, "Show the graph of dependencies")
	flag.StringVar(&hostfileName, "hostfile", "hostfile.cfg", "File listing host running the listener")
	flag.StringVar(&makefileName, "makefile", "Makefile", "The Makefile")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	log.Println("Parsing the Makefile...")
	head, err := parser.Parse(makefileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Done")

	if showGraph {
		log.Print("Graph:\n", head)
	}

	log.Println("Parsing the hostfile...")
	hosts := config.Parse(hostfileName)
	log.Println("Done")

	todo := make(chan *parser.Task)

	for _, host := range hosts {
		go run(host, todo, verbose)
	}

	for !walk(head, todo) {
	}

}
