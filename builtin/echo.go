package builtin

import (
	"fmt"
	"strings"
)

func init() {
	Register("echo", Echo, "Display line of text")
}

func Echo(args []string) {
	var output []string
	newline := true
	
	// Parse flags
	i := 0
	for i < len(args) {
		if args[i] == "-n" {
			newline = false
			i++
		} else if args[i] == "-e" {
			// Enable interpretation of backslash escapes (we'll implement basic ones)
			i++
			// For now, we'll just skip this flag and handle escapes by default
		} else {
			break
		}
	}
	
	// Collect remaining arguments as text to echo
	for j := i; j < len(args); j++ {
		output = append(output, args[j])
	}
	
	text := strings.Join(output, " ")
	
	// Handle basic escape sequences
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\t", "\t")
	text = strings.ReplaceAll(text, "\\r", "\r")
	text = strings.ReplaceAll(text, "\\\\", "\\")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "\\'", "'")
	
	if newline {
		fmt.Println(text)
	} else {
		fmt.Print(text)
	}
}