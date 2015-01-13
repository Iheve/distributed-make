package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Task struct {
	Target   string
	Deps     []string
	Cmds     []string
	Sons     []*Task
	Affected bool
	Done     bool
}

func readTarget(l string) (target string, deps []string) {
	//Get the target
	if !strings.Contains(l, ":") {
		log.Fatal("Invalid line : can't find separator ':' in line : ", l)
	}
	c := strings.SplitN(l, ":", 2)
	target = strings.TrimSpace(c[0])
	//Get the dependencies
	deps = strings.Split(strings.TrimSpace(c[1]), " ")
	return
}

func readCmd(l string) (cmds []string) {
	cmds = nil
	for _, c := range strings.Split(l, ";") {
		cmds = append(cmds, strings.TrimSpace(c))
	}
	return
}

func newTask(target string, deps []string, cmds []string) *Task {
	task := new(Task)
	task.Target = target
	task.Deps = deps
	task.Cmds = cmds
	task.Affected = false
	task.Affected = false
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
	var cmds []string = nil
	targetSet := false
	first := true

	for scanner.Scan() {
		if len(scanner.Text()) == 0 || strings.HasPrefix(scanner.Text(), "#") {
			//Skip empty lines and comments
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
			if targetSet {
				tasks[target] = newTask(target, deps, cmds)
				if first {
					head = tasks[target]
					first = false
				}
				cmds = nil
			}
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

func walk(t *Task, d int, buffer *bytes.Buffer) {
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString(t.Target + "\n")

	for _, c := range t.Cmds {
		for i := 0; i < d; i++ {
			buffer.WriteString("\t")
		}
		buffer.WriteString(c + "\n")
	}

	for _, s := range t.Sons {
		if s != nil {
			walk(s, d+1, buffer)
		}
	}
}

func (t *Task) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	walk(t, 0, &buffer)
	return buffer.String()
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

	fmt.Print(head)
}
