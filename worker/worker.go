package worker

import (
	"fmt"
	"os/exec"
)

type Args struct {
	Command string
}

type Response struct {
	Output string
}

type Worker int

func (t *Worker) Output(args *Args, response *Response) error {
	out, err := exec.Command(args.Command).Output()
	response.Output = fmt.Sprintf("%s", out)
	if err == nil {
		fmt.Printf("Command executed successfully. Output:\n%s", out)
	} else {
		fmt.Printf("Command failed with error: %v", err)
	}
	return err
}
