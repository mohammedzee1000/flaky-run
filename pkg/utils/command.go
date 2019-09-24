package utils

import (
	"os/exec"
	"strings"
)

//RunCommand runs a command
func RunCommand(name string, args ...string) (string, error) {
	c := exec.Command(name, args...)
	o, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(o), "\n"), err
}
