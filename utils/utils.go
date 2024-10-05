package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/SanjaySinghRajpoot/realCode/utils/redis"
)

// let's try to optimize this part if possible
func CompileCodePython(code string, userid uint) (string, error) {

	// apply a function to get the code from the redis cache
	codeOutput, _ := redis.GetCode(code)

	print("Code Output: ", codeOutput)

	// fmt.Printf("Code Output: %s", codeOutput)

	// if codeOutput != "" {
	// 	return codeOutput, nil
	// }

	cmd := exec.Command("bash", "-c", "echo '"+code+"' | python3 -c 'import sys; exec(sys.stdin.read())'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}

	// here we will implement redis cache for code
	msg, err := redis.SetCode(code, out.String())
	if err != nil {
		fmt.Printf("Failed to Set the Redis Cache CRON: %s", msg)

		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}

	return out.String(), nil
}

func CompileCodeGo(code string, userid uint) (string, error) {

	cmd := exec.Command("bash", "-c", "echo '"+code+"' > temp.go && go run temp.go && rm temp.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// here we will implement redis cache for code
	msg, err := redis.SetCode(code, out.String())
	if err != nil {
		fmt.Printf("Failed to Set the Redis Cache CRON: %s", msg)

		errRes := fmt.Sprintf("error: %s", err.Error())
		return errRes, err
	}

	return out.String(), nil
}
