package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	// "os/signal"
	"github.com/g3rzi/ctenterd/builtin"
	"github.com/g3rzi/ctenterd/pkg/color"
    "os/exec"      
    "path/filepath" 
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var (
	version   = "v0.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-V" {
			fmt.Printf("ctenterd %s\n", version)
			return
		}
		line := strings.Join(os.Args[1:], " ")
		if err := runProgram(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	runREPL()
}

// func init() {
// 	// Swallow SIGINT so Ctrl+C doesn't kill the whole shell.
// 	sigc := make(chan os.Signal, 1)
// 	signal.Notify(sigc, os.Interrupt)
// 	go func() {
// 		for range sigc {
// 			// do nothing: this consumes SIGINT and keeps the shell alive
// 		}
// 	}()
// }

// Read–Eval–Print Loop
func runREPL() {
	fmt.Printf("ctenterd %s - Container Shell Agent\n", version)
	fmt.Println("Type 'help' for available commands, 'exit' to quit")

	sc := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(color.Prompt("ctenterd> "))
		if !sc.Scan() {
			break
		}

		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			break
		}

		if line == "help" {
			printHelp(os.Stdout)
			continue
		}

		if err := runProgram(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}
}

func runProgram(src string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(src), "")
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	r, err := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.ExecHandler(execBuiltins), // signature: func(ctx context.Context, args []string) error
	)
	if err != nil {
		return fmt.Errorf("interp init: %w", err)
	}

	return r.Run(context.Background(), file)
}

// execBuiltins runs our built-ins. Return nil on success; error to surface failure.
// func execBuiltins(ctx context.Context, args []string) error {
// 	if len(args) == 0 {
// 		return nil
// 	}

// 	// built-in help: list or per-command
// 	if args[0] == "help" {
// 		if len(args) == 1 {
// 			// List all commands with their help text
// 			names := builtin.List()
// 			if len(names) == 0 {
// 				fmt.Println("No commands registered.")
// 				return nil
// 			}

// 			// Find padding width based on the longest command name
// 			max := 0
// 			for _, name := range names {
// 				if len(name) > max {
// 					max = len(name)
// 				}
// 			}

// 			fmt.Println("Available commands:")
// 			for _, name := range names {
// 				if help, ok := builtin.Help(name); ok {
// 					fmt.Printf("  %-*s  %s\n", max, name, help)
// 				}
// 			}
// 			return nil
// 		}

// 		// Help for specific command
// 		if help, ok := builtin.Help(args[1]); ok {
// 			fmt.Printf("%s: %s\n", args[1], help)
// 			return nil
// 		}
// 		return fmt.Errorf("%s: no help available", args[1])
// 	}

// 	if fn, ok := builtin.Lookup(args[0]); ok {
// 		fn(args[1:])
// 		return nil
// 	}

// 	// fallback to toybox/PATH/etc. here...
// 	return fmt.Errorf("%s: not found (no builtin and no toolbox)", args[0])
// }

func execBuiltins(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	// built-in help stays as you wrote it
	if args[0] == "help" {
		if len(args) == 1 {
			names := builtin.List()
			if len(names) == 0 {
				fmt.Println("No commands registered.")
				return nil
			}
			max := 0
			for _, name := range names {
				if len(name) > max {
					max = len(name)
				}
			}
			fmt.Println("Available commands:")
			for _, name := range names {
				if help, ok := builtin.Help(name); ok {
					fmt.Printf("  %-*s  %s\n", max, name, help)
				}
			}
			return nil
		}
		if help, ok := builtin.Help(args[1]); ok {
			fmt.Printf("%s: %s\n", args[1], help)
			return nil
		}
		return fmt.Errorf("%s: no help available", args[1])
	}

	// 1) Builtins
	if fn, ok := builtin.Lookup(args[0]); ok {
		fn(args[1:])
		return nil
	}

	// 2) External: path or PATH
	name := args[0]
	argv := args[1:]

	// Looks like an explicit path? (absolute or contains '/')
	if strings.ContainsRune(name, os.PathSeparator) {
		return runExternal(ctx, name, argv)
	}

	// Otherwise, search PATH
	lp, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s: not found", name)
	}
	return runExternal(ctx, lp, argv)
}

func runExternal(ctx context.Context, path string, argv []string) error {
	// Optional: resolve to absolute for nicer errors
	if !filepath.IsAbs(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}

	// Validate exists & not a dir (gives clearer errors than exec alone)
	if st, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s: %v", path, err)
	} else if st.IsDir() {
		return fmt.Errorf("%s: is a directory", path)
	}

	cmd := exec.CommandContext(ctx, path, argv...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // inherit env

	if err := cmd.Run(); err != nil {
		// Common “gotchas” to print a helpful hint
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no such file or directory"):
			// Missing interpreter (script without shebang) or missing dynamic loader
			return fmt.Errorf("%s: exec failed (missing interpreter/shebang or dynamic loader)", path)
		case strings.Contains(msg, "permission denied"):
			return fmt.Errorf("%s: permission denied (chmod +x ?)", path)
		case strings.Contains(msg, "exec format error"):
			return fmt.Errorf("%s: exec format error (wrong arch?)", path)
		default:
			return fmt.Errorf("%s: %v", path, err)
		}
	}
	return nil
}


func printHelp(w io.Writer) {
	names := builtin.List()
	if len(names) == 0 {
		fmt.Fprintln(w, "No commands registered.")
		return
	}

	// Find padding width based on the longest command name
	max := 0
	for _, name := range names {
		if len(name) > max {
			max = len(name)
		}
	}

	fmt.Fprintln(w, "Available commands:")
	for _, name := range names {
		if help, ok := builtin.Help(name); ok {
			fmt.Fprintf(w, "  %-*s  %s\n", max, name, help)
		}
	}
	fmt.Fprintln(w, "\nTip: type 'help <command>' for specific help or 'help' to see this again.")
}