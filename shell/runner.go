// shell/runner.go
package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/g3rzi/ctenterd/builtin"
)

func ExecLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	args := fields(line)
	name := args[0]
	argv := args[1:]

	// 1) Builtins
	if fn, ok := builtin.Lookup(name); ok {
		fn(argv)
		return
	}

	// 2) If command looks like a path (absolute or contains '/'), run it directly
	if looksLikePath(name) {
		runExec(name, argv)
		return
	}

	// 3) Otherwise search $PATH and run
	if lp, err := exec.LookPath(name); err == nil {
		runExec(lp, argv)
		return
	}

	fmt.Fprintf(os.Stderr, "%s: not found\n", name)
}

func runExec(path string, argv []string) {
	if !filepath.IsAbs(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}
	st, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
		return
	}
	if st.IsDir() {
		fmt.Fprintf(os.Stderr, "%s: is a directory\n", path)
		return
	}

	cmd := exec.Command(path, argv...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Fprintf(os.Stderr, "%s: exec format or missing interpreter (check shebang/ld-linux)\n", path)
			return
		}
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
	}
}

func looksLikePath(s string) bool {
	return filepath.IsAbs(s) || strings.Contains(s, "/")
}

func fields(s string) []string { return strings.Fields(s) }
