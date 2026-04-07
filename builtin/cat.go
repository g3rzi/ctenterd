package builtin

import (
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	Register("cat", Cat, "Print file contents")
}

func Cat(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "cat: missing file argument\n")
		return
	}

	
	
	for _, filename := range args {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
			continue
		}
		
		fmt.Print(string(data))
	}
}
