package builtin

import (
	"fmt"
	"os"
)

func init() {
	Register("mkdir", Mkdir, "Create directories")
}

// mkdir - Create directories
func Mkdir(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "mkdir: missing operand\n")
		return
	}

	parents := false
	var dirs []string

	// Parse arguments
	for _, arg := range args {
		if arg == "-p" {
			parents = true
		} else {
			dirs = append(dirs, arg)
		}
	}

	for _, dir := range dirs {
		var err error
		if parents {
			err = os.MkdirAll(dir, 0755)
		} else {
			err = os.Mkdir(dir, 0755)
		}
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: cannot create directory '%s': %v\n", dir, err)
		}
	}
}