package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Target string
	Deps   []string
	Cmd    *exec.Cmd
	Sons   []*Task
}

func line1(l string) (target string, deps []string) {
	//Get the target
	if !strings.Contains(l, ":") {
		log.Fatal("Invalid line : can't find separator ':'")
	}
	c := strings.SplitN(l, ":", 2)
	target = strings.TrimSpace(c[0])
	//Get the dependencies
	deps = strings.Split(strings.TrimSpace(c[1]), " ")
	return
}

func line2(l string, target string, deps []string) (task *Task) {
	//Build the command
	c := strings.TrimSpace(l)
	args := strings.Split(c, " ")
	cmd := exec.Command(args[0])
	cmd.Args = args
	//Build the task
	task = new(Task)
	task.Target = target
	task.Deps = deps
	task.Cmd = cmd
	return
}

func linkTasks(tasks map[string]*Task) {
	for _, t := range tasks {
		for _, d := range t.Deps {
			t.Sons = append(t.Sons, tasks[d])
		}
	}
}

func Parse(filename string) (head *Task, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	tasks := make(map[string]*Task)
	var target string
	var deps []string

	for first := true; scanner.Scan(); {
		if len(scanner.Text()) == 0 {
			//Skip empty lines
			continue
		}
		if strings.HasPrefix(scanner.Text(), "\t") {
			tasks[target] = line2(scanner.Text(), target, deps)
			if first {
				head = tasks[target]
				first = false
			}
		} else {
			target, deps = line1(scanner.Text())
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

	linkTasks(tasks)

	return
}

func walk(t *Task, d int) {
	for i := 0; i < d; i++ {
		fmt.Print("\t")
	}
	fmt.Println(t.Target)
	/*
		for i := 0; i < d; i++ {
			fmt.Print("\t")
		}
		fmt.Println(t.Cmd.Args)
	*/
	for _, s := range t.Sons {
		if s != nil {
			walk(s, d+1)
		}
	}
}

func main() {
	// Check if there is an argument
	if len(os.Args) != 2 {
		path = "Makefile"
	} else {
		path := os.Args[1]
	}

	head, err := Parse(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Println(head)
	walk(head, 0)
}
