package builtin

import (
	"fmt"
	"os"
	"sort"
)

func init() {
	Register("env", Env, "Show environment variables")
}

func Env(args []string) {
	env := os.Environ()
	sort.Strings(env)
	
	for _, e := range env {
		fmt.Println(e)
	}
}
