package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/SanjaySinghRajpoot/realCode/utils/redis"
)

// let's try to optimize this part if possible
func CompileCodePython(code string, userid uint) (string, error) {

	// here we will implement redis cache for code
	msg, err := redis.SetCode(code, userid)
	if err != nil {
		fmt.Printf("Failed to Set the Redis Cache CRON: %s", msg)

		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}

	cmd := exec.Command("bash", "-c", "echo '"+code+"' | python3 -c 'import sys; exec(sys.stdin.read())'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}
	return out.String(), nil
}

func CompileCodeGo(code string, userid uint) (string, error) {

	msg, err := redis.SetCode(code, userid)
	if err != nil {
		fmt.Printf("Failed to Set the Redis Cache CRON: %s", msg)

		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}

	cmd := exec.Command("bash", "-c", "echo '"+code+"' > temp.go && go run temp.go && rm temp.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
