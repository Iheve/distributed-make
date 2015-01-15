package main

import (
	"flag"
	"fmt"
	"github.com/Iheve/distributed-make/config"
	"github.com/Iheve/distributed-make/parser"
	"github.com/Iheve/distributed-make/worker"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"
)

type job struct {
	name      string
	startTime time.Time
}

var running, finished chan job
var pretty bool

func run(host string, todo chan *parser.Task, verbose, showTimes bool) {
	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		if !pretty {
			log.Println("Can not contact", host, err)
		}
		return
	}
	for {
		t := <-todo
		id := job{fmt.Sprintf("%v:%v", host, t.Target), time.Now()}
		if pretty {
			running <- id
		} else {
			log.Println(host, "builds", t.Target)
		}
		now := time.Now()
		var response worker.Response
		args := new(worker.Args)
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
			if pretty {
				pretty = false
				termbox.Close()
			}
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

		if showTimes {
			log.Printf("%s has built %s in %v", host, t.Target, time.Since(now))
		}
		t.Done = true
		if verbose {
			log.Println(host, "Command done, outputs:")
			for _, s := range response.Output {
				log.Println("\n", s)
			}
		}
		if pretty {
			finished <- id
		}
	}
}

func walk(t *parser.Task, todo chan *parser.Task) (bool, time.Time) {
	if t.Done {
		fileInfo, _ := os.Stat(t.Target)
		return true, fileInfo.ModTime()
	}

	if t.Affected {
		return false, time.Unix(0, 0)
	}

	res := true
	mostRecentCreationDate := time.Unix(0, 0)
	for _, s := range t.Sons {
		if s != nil {
			done, creationDateSon := walk(s, todo)
			res = done && res
			if res {
				if mostRecentCreationDate.Before(creationDateSon) {
					mostRecentCreationDate = creationDateSon
				}
			}
		}
	}

	if res {
		if fileInfo, err := os.Stat(t.Target); err == nil && fileInfo.ModTime().After(mostRecentCreationDate) {
			t.Done = true
			return true, fileInfo.ModTime()
		}
		t.Affected = true
		todo <- t
	}

	return false, time.Unix(0, 0)
}

func writeString(x, y int, j job) {
	s := fmt.Sprintf("%s %v", j.name, time.Since(j.startTime))
	for i, r := range s {
		termbox.SetCell(x+i, y, r, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func writeList(x, y int, l []job) {
	for i, s := range l {
		writeString(x, y+i, s)
	}
}

func events(c chan int) {
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEsc {
			pretty = false
			c <- 0
		}

	}
}

func display() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	quit := make(chan int)

	go events(quit)

	var l []job = nil
	for {
		select {
		case j := <-running:
			l = append(l, j)
		case j := <-finished:
			for i := range l {
				if l[i] == j {
					l = append(l[:i], l[i+1:]...)
					break
				}
			}
		case <-quit:
			termbox.Close()
			log.Println("Exiting pretty mode")
			return
		default:
			time.Sleep(1e8)
		}
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		writeList(0, 0, l)
		termbox.Flush()
	}
}

func first(b bool, t time.Time) bool {
	return b
}

func main() {
	var help, verbose, showGraph, showTimes bool
	var hostfileName, makefileName string
	var nbThread int
	flag.BoolVar(&help, "help", false, "Display this helper message")
	flag.BoolVar(&verbose, "verbose", false, "Show outputs of commands")
	flag.BoolVar(&pretty, "pretty", false, "Display a pretty output")
	flag.BoolVar(&showTimes, "showtimes", false, "Show in how much time the target has been built")
	flag.BoolVar(&showGraph, "showgraph", false, "Show the graph of dependencies")
	flag.StringVar(&hostfileName, "hostfile", "hostfile.cfg", "File listing host running the listener")
	flag.StringVar(&makefileName, "makefile", "Makefile", "The Makefile")
	flag.IntVar(&nbThread, "nbthread", 1, "Number of thread per worker")
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
	running = make(chan job, 1000)
	finished = make(chan job, 1000)

	for i := 0; i < nbThread; i++ {
		for _, host := range hosts {
			go run(host, todo, verbose, showTimes)
		}
	}

	if pretty {
		go display()
	}

	for !first(walk(head, todo)) {
	}

	if pretty {
		termbox.Close()
	}

}
