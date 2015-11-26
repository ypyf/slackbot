package plugin

import (
	"os/exec"
)

func ExecShell(bin string, args...string) (string, error) {
	cmd := exec.Command(bin, args...)
	stdout, err := cmd.Output()
	return string(stdout), err
}
