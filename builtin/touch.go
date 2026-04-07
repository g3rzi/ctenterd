package builtin

import (
	"fmt"
	"os"
	"time"
)

func init() {
	Register("touch", Touch, "Create empty files or update timestamps")
}



// touch - Create empty files or update timestamps
func Touch(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "touch: missing file operand\n")
		return
	}

	now := time.Now()

	for _, filename := range args {
		// Try to update timestamp if file exists
		if err := os.Chtimes(filename, now, now); err != nil {
			// File doesn't exist, create it
			if os.IsNotExist(err) {
				file, createErr := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
				if createErr != nil {
					fmt.Fprintf(os.Stderr, "touch: cannot touch '%s': %v\n", filename, createErr)
					continue
				}
				file.Close()
			} else {
				fmt.Fprintf(os.Stderr, "touch: cannot touch '%s': %v\n", filename, err)
			}
		}
	}
}