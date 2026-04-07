package builtin

import (
	"fmt"
	"os"
)

func init() {
	Register("pwd", Pwd, "Print working directory")
}

func Pwd(args []string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return
	}
	
	fmt.Println(pwd)
}
