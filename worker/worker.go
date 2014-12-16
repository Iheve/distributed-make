package worker

import (
	"fmt"
	"os/exec"
)

type Args struct {
	Command   string
	Arguments []string
}

type Response struct {
	Output string
}

type Worker int

func (t *Worker) Output(args *Args, response *Response) error {
	cmd := exec.Command(args.Command)
	cmd.Args = args.Arguments
	out, err := cmd.Output()
	response.Output = fmt.Sprintf("%s", out)
	if err == nil {
		fmt.Printf("Command executed successfully. Output:\n%s", out)
	} else {
		fmt.Printf("Command failed with error: %v", err)
	}
	return err
}
