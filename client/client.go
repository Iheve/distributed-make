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

var pretty bool

var runningJobs, doneJobs []job
var running, finished chan job = make(chan job, 1000), make(chan job, 1000)
var hosts, failedHosts []string
var addHost, rmHost chan string = make(chan string, 1000), make(chan string, 1000)

func run(host string, todo chan *parser.Task, verbose, showTimes bool) {
	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		rmHost <- host
		if !pretty {
			log.Println("Can not contact", host, err)
		}
		return
	}
	for {
		t := <-todo
		id := job{fmt.Sprintf("%v:%v", host, t.Target), time.Now()}
		running <- id
		if !pretty {
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
				if pretty {
					pretty = false
					termbox.Close()
				}
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
				rmHost <- host
				if !pretty {
					log.Println("Contact lost with ", host)
					log.Println(t.Target, "will be rebuilt.")
					log.Println(host, "will not receive job anymore.")
				}
				todo <- t
				return
			}
			if pretty {
				pretty = false
				termbox.Close()
			}
			log.Fatal(host, " failed to build target ", t.Target, ":", err)
		}
		//Unpack target
		err = ioutil.WriteFile(response.Target.Name, response.Target.Content, response.Target.Mode)
		if err != nil {
			if pretty {
				pretty = false
				termbox.Close()
			}
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

func events() {
	for pretty {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEsc {
			pretty = false
		}
	}
}

func updateStatus() {
	for {
		select {
		case host := <-addHost:
			hosts = append(hosts, host)
		case host := <-rmHost:
			failedHosts = append(failedHosts, host)
			for i, h := range hosts {
				if h == host {
					hosts = append(hosts[:i], hosts[i+1:]...)
					break
				}
			}
			if len(hosts) == 0 {
				if pretty {
					pretty = false
					termbox.Close()
				}
				log.Println("No more worker available")
				os.Exit(2)
			}
		case job := <-running:
			runningJobs = append(runningJobs, job)
		case job := <-finished:
			doneJobs = append(doneJobs, job)
			for i, j := range runningJobs {
				if j == job {
					runningJobs = append(runningJobs[:i], runningJobs[i+1:]...)
					break
				}
			}
		}
	}
}

func writeString(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func writeList(x int, title string, l ...string) {
	writeString(x, 0, title)
	for i, s := range l {
		writeString(x, i+1, s)
	}
}

func display() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	go events()

	for pretty {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		var l []string
		for _, j := range runningJobs {
			timeSince := time.Since(j.startTime)
			var divider float32
			var unit string
			if timeSince > time.Second {
				divider = float32(time.Second)
				unit = "s"
			} else {
				divider = float32(time.Millisecond)
				unit = "ms"
			}
			var duration = float32(timeSince) / divider
			l = append(l, fmt.Sprintf("%s:%.3f%s", j.name, duration, unit))
		}
		writeList(0, "Running jobs", l...)
		l = nil
		for _, j := range doneJobs {
			l = append(l, fmt.Sprintf("%s", j.name))
		}
		writeList(45, "Done jobs", l...)
		writeList(75, "Hosts", hosts...)
		writeList(95, "Failed hosts", failedHosts...)
		termbox.Flush()
	}
}

func first(b bool, t time.Time) bool {
	return b
}

func main() {
	var help, verbose, showGraph, showTimes bool
	var hostfileName, makefileName, target string
	var nbThread int
	flag.BoolVar(&help, "help", false, "Display this helper message")
	flag.BoolVar(&verbose, "verbose", false, "Show outputs of commands")
	flag.BoolVar(&pretty, "pretty", false, "Display a pretty output")
	flag.BoolVar(&showTimes, "showtimes", false, "Show in how much time the target has been built")
	flag.BoolVar(&showGraph, "showgraph", false, "Show the graph of dependencies")
	flag.StringVar(&hostfileName, "hostfile", "hostfile.cfg", "File listing host running the listener")
	flag.StringVar(&makefileName, "makefile", "Makefile", "The Makefile")
	flag.StringVar(&target, "target", "", "Target to execute")
	flag.IntVar(&nbThread, "nbthread", 1, "Number of thread per worker")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	log.Println("Parsing the Makefile...")
	head, err := parser.Parse(makefileName, target)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Done")

	if showGraph {
		log.Print("Graph:\n", head)
	}

	log.Println("Parsing the hostfile...")
	myHosts := config.Parse(hostfileName)
	log.Println("Done")

	if len(myHosts) == 0 {
		log.Println("No worker available, please add listeners to hostfile")
		os.Exit(2)
	}

	go updateStatus()

	todo := make(chan *parser.Task)

	for i := 0; i < nbThread; i++ {
		for _, host := range myHosts {
			addHost <- host
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

	fmt.Println(len(hosts))
}
