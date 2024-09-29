package pkg

import (
	"fmt"
	"os/exec"
)

func Xk(command string, options map[string]string) error {
	cmd := exec.Command("xk", command)

	var args []string
	for key, value := range options {
		args = append(args, fmt.Sprintf("-%s", key))
		args = append(args, value)
	}

	cmd.Args = append(cmd.Args, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing script: %v, output: %s", err, output)
	}

	fmt.Println(string(output))

	return nil
}
