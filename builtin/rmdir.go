package builtin

import (
	"fmt"
	"os"
)

func init() {
	Register("rmdir", Rmdir, "Remove empty directories")
}

// rmdir - Remove empty directories
func Rmdir(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "rmdir: missing operand\n")
		return
	}

	for _, dir := range args {
		if err := os.Remove(dir); err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: failed to remove '%s': %v\n", dir, err)
		}
	}
}