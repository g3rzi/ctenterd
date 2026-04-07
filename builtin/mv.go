package builtin

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	Register("mv", Move, "Move/rename files or directories")
}


// mv - Move/rename files or directories
func Move(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "mv: missing file arguments\n")
		fmt.Fprintf(os.Stderr, "Usage: mv SOURCE DEST\n")
		return
	}

	if len(args) == 2 {
		// Simple rename/move
		if err := os.Rename(args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "mv: cannot move '%s' to '%s': %v\n", args[0], args[1], err)
		}
		return
	}

	// Multiple sources - destination must be a directory
	dest := args[len(args)-1]
	sources := args[:len(args)-1]

	destInfo, err := os.Stat(dest)
	if err != nil || !destInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "mv: target '%s' is not a directory\n", dest)
		return
	}

	for _, src := range sources {
		basename := filepath.Base(src)
		destPath := filepath.Join(dest, basename)
		if err := os.Rename(src, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "mv: cannot move '%s' to '%s': %v\n", src, destPath, err)
		}
	}
}