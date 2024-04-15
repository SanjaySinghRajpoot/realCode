package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// let's try to optimize this part if possible
func CompileCodePython(code string) (string, error) {
	cmd := exec.Command("bash", "-c", "echo '"+code+"' | python3 -c 'import sys; exec(sys.stdin.read())'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}
	return out.String(), nil
}

func CompileCodeGo(code string) (string, error) {
	cmd := exec.Command("bash", "-c", "echo '"+code+"' > temp.go && go run temp.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
