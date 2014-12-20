package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

type File struct {
	Name    string
	Content []byte
}

type Args struct {
	Target string
	Cmds   [][]string
	Deps   []File
}

type Response struct {
	Output []string
	Target File
}

type Worker int

func (t *Worker) Output(args *Args, response *Response) error {
	log.Println("New rpc")
	//Create temp dir
	dir, err := ioutil.TempDir("", "dmake")
	if err != nil {
		log.Println("Can not create temp dir: ", err)
		return err
	}
	dir = dir + "/"
	log.Println("Using temp dir: ", dir)
	//Unpack dependencies
	for _, f := range args.Deps {
		if f.Name == "" {
			continue
		}
		err := ioutil.WriteFile(dir+f.Name, f.Content, 0777)
		if err != nil {
			log.Println("Can not create file: ", f.Name, " : ", err)
			return err
		}
	}
	//Run commands
	for _, cmd := range args.Cmds {
		log.Println("Executing cmd :", cmd)
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Dir = dir
		out, err := c.Output()
		response.Output = append(response.Output, fmt.Sprintf("%s", out))
		if err == nil {
			log.Printf("Command executed successfully. Output:\n%s", out)
		} else {
			log.Printf("Command failed with error: %v", err)
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

	return nil
}
