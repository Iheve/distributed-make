package worker

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type File struct {
	Name    string
	Mode    os.FileMode
	Content []byte
}

type Args struct {
	Target string
	Cmds   []string
	Deps   []File
}

type Response struct {
	Output []string
	Target File
}

func execute(cmd, dir string, env []string) (outCmd []byte, err error) {
	if strings.ContainsAny(cmd, ";><$`&") {
		c := exec.Command("bash", "-c", cmd)
		c.Dir = dir
		c.Env = env
		return c.CombinedOutput()
	}
	var output bytes.Buffer
	var cmdSplit = strings.Split(cmd, "|")
	var cmds []*exec.Cmd
	var pipes = make([]*io.PipeWriter, len(cmdSplit)-1)
	var i int

	for _, oneCmd := range cmdSplit {
		oneCmd := strings.TrimSpace(oneCmd)
		args := strings.Split(oneCmd, " ")
		c := exec.Command(args[0], args[1:]...)
		c.Dir = dir
		c.Env = env
		cmds = append(cmds, c)
	}

	for i = 0; i < len(cmds)-1; i++ {
		in, out := io.Pipe()
		cmds[i].Stdout = out
		cmds[i+1].Stdin = in
		pipes[i] = out
	}
	cmds[i].Stdout = &output
	cmds[i].Stderr = &output

	if err := call(cmds, pipes); err != nil {
		log.Printf("Command failed with error: %v", err)
	}
	return output.Bytes(), err
}

func call(cmds []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if cmds[0].Process == nil {
		if err = cmds[0].Start(); err != nil {
			return err
		}
	}
	if len(cmds) > 1 {
		if err = cmds[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = call(cmds[1:], pipes[1:])
			}
		}()
	}
	return cmds[0].Wait()
}

type Worker int

func (t *Worker) Output(args *Args, response *Response) error {
	log.Println("Building target: ", args.Target)
	//Create temp dir
	dir, err := ioutil.TempDir("", "dmake")
	if err != nil {
		log.Println("Can not create temp dir: ", err)
		return err
	}
	dir = dir + "/"

	//Calculate new environment (set PWD)
	env := os.Environ()
	for i := range env {
		if strings.HasPrefix(env[i], "PWD") {
			env[i] = "PWD=" + dir
			break
		}
	}

	//Unpack dependencies
	for _, f := range args.Deps {
		if f.Name == "" {
			continue
		}
		err := ioutil.WriteFile(dir+f.Name, f.Content, f.Mode)
		if err != nil {
			log.Println("Can not create file: ", f.Name, " : ", err)
			return err
		}
	}
	//Run commands
	for _, cmd := range args.Cmds {
		out, err := execute(cmd, dir, env)
		response.Output = append(response.Output, fmt.Sprintf("%s", out))
		log.Println(string(out))
		if err == nil {
		} else {
			log.Println("Command failed with error ", err, " output:")
			log.Println(string(out))
			return err
		}
	}
	//Pack target
	response.Target.Name = args.Target
	response.Target.Content, err = ioutil.ReadFile(dir + args.Target)
	if err != nil {
		log.Println("Cant read file: ", args.Target, " : ", err)
		return err
	}
	info, _ := os.Stat(dir + args.Target)
	response.Target.Mode = info.Mode()

	//Rm temp dir
	os.RemoveAll(dir)

	return nil
}
