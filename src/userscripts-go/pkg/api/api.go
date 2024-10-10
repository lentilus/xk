package api

import (
	"fmt"
	"os/exec"
	"strings"
)

func Xk(command string, options map[string]string) ([]string, error) {
	cmd := exec.Command("xk", command)

	var args []string
	for key, value := range options {
		args = append(args, fmt.Sprintf("-%s", key))
		args = append(args, value)
	}

	cmd.Args = append(cmd.Args, args...)

	output, err := cmd.Output()
	if err != nil {
		return []string{}, fmt.Errorf("error executing script: %v, output: %s", err, output)
	}

	return strings.Split(string(output), "\n"), nil
}
