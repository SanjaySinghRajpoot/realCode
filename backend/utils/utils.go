package utils

import (
	"bytes"
	"os/exec"
)

// let's try to optimize this part if possible
func CompileCode(code string) (string, error) {
	cmd := exec.Command("bash", "-c", "echo '"+code+"' | python3 -c 'import sys; exec(sys.stdin.read())'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
