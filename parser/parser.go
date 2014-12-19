package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Target string
	Deps   []string
	Cmds   []*exec.Cmd
	Sons   []*Task
}

func readTarget(l string) (target string, deps []string) {
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

func readCmd(l string) (cmds []*exec.Cmd) {
	cmds = nil
	for _, c := range strings.Split(l, ";") {
		c := strings.TrimSpace(c)
		args := strings.Split(c, " ")
		cmd := exec.Command(args[0], args...)
		cmds = append(cmds, cmd)
	}
	return
}

func newTask(target string, deps []string, cmds []*exec.Cmd) *Task {
	task := new(Task)
	task.Target = target
	task.Deps = deps
	task.Cmds = cmds
	return task
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
	var cmds []*exec.Cmd = nil
	targetSet := false
	first := true

	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			//Skip empty lines
			if targetSet {
				tasks[target] = newTask(target, deps, cmds)
				if first {
					head = tasks[target]
					first = false
				}
				cmds = nil
			}
			targetSet = false
			continue
		}
		if strings.HasPrefix(scanner.Text(), "\t") {
			if !targetSet {
				err = errors.New("Parser : target must be set before the command line")
				log.Fatal(err)
				return
			}
			cmds = append(cmds, readCmd(scanner.Text())...)
		} else {
			target, deps = readTarget(scanner.Text())
			targetSet = true
		}
	}

	if targetSet {
		tasks[target] = newTask(target, deps, cmds)
		if first {
			head = tasks[target]
			first = false
		}
		cmds = nil
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

	for _, c := range t.Cmds {
		for i := 0; i < d; i++ {
			fmt.Print("\t")
		}
		fmt.Println(c.Args)
	}

	for _, s := range t.Sons {
		if s != nil {
			walk(s, d+1)
		}
	}
}

func main() {
	// Check if there is an argument
	var path string
	if len(os.Args) != 2 {
		path = "Makefile"
	} else {
		path = os.Args[1]
	}

	head, err := Parse(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	walk(head, 0)
}
