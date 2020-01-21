package tools

import (
	"fmt"
	"os/exec"
)

func GitPush() {
	command := `./tools/git_push.sh .`
	cmd := exec.Command("/bin/bash", "-c", command)

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}