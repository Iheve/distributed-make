package main

import (
	"flag"
	"github.com/Iheve/distributed-make/config"
	"github.com/Iheve/distributed-make/parser"
	"github.com/Iheve/distributed-make/worker"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

func run(client *rpc.Client, name string, todo chan *parser.Task, verbose bool) {
	for {
		t := <-todo
		var response worker.Response
		args := new(worker.Args)
		log.Println(name, "builds", t.Target)
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
		//Synchronous call
		err := client.Call("Worker.Output", args, &response)
		if err != nil {
			log.Fatal(name, "RPC call error:", err)
		}
		//Unpack target
		err = ioutil.WriteFile(response.Target.Name, response.Target.Content, response.Target.Mode)
		if err != nil {
			log.Fatal("Can not create file: ", response.Target.Name, " : ", err)
		}

		t.Done = true
		if verbose {
			log.Println(name, "Command done, outputs:")
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

func findTasks(chanDone chan int, done *bool, head *parser.Task, todo chan *parser.Task) {
	for !walk(head, todo) {
	}
	chanDone <- 1
	*done = true
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
	chanDone := make(chan int)

	done := false
	go findTasks(chanDone, &done, head, todo)

	for _, host := range hosts {
		if done {
			break
		}
		client, err := rpc.DialHTTP("tcp", host)
		if err != nil {
			log.Println("Can not contact", host, err)
			continue
		}
		go run(client, host, todo, verbose)
	}

	<-chanDone

}
