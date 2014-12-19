package worker

import (
	"fmt"
	"os/exec"
)

type Args struct {
	Cmds [][]string
}

type Response struct {
	Output []string
}

type Worker int

func (t *Worker) Output(args *Args, response *Response) error {
	fmt.Println("New rpc")
	for _, cmd := range args.Cmds {
		c := exec.Command(cmd[0], cmd[1:]...)
		out, err := c.Output()
		response.Output = append(response.Output, fmt.Sprintf("%s", out))
		if err == nil {
			fmt.Printf("Command executed successfully. Output:\n%s", out)
		} else {
			fmt.Printf("Command failed with error: %v", err)
			return err
		}
	}
	return nil
}
