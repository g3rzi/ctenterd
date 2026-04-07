package builtin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Register("cd", CD, "Change the shell working directory")
}

func CD(args []string) {
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "cd: too many arguments")
		return
	}

	// figure target
	var target string
	if len(args) == 0 || args[0] == "~" {
		home := getenvHome()
		if home == "" {
			fmt.Fprintln(os.Stderr, "cd: HOME not set")
			return
		}
		target = home
	} else {
		arg := args[0]
		switch arg {
		case "-":
			old := os.Getenv("OLDPWD")
			if old == "" {
				fmt.Fprintln(os.Stderr, "cd: OLDPWD not set")
				return
			}
			target = old
			// bash prints the path when using "cd -"
			fmt.Println(target)
		default:
			target = expandTilde(arg)
		}
	}

	// handle CDPATH only for relative paths (like bash)
	cdpathHit := false
	if !filepath.IsAbs(target) && target != "." && target != ".." {
		if cdpath := os.Getenv("CDPATH"); cdpath != "" {
			for _, base := range splitList(cdpath) {
				// empty entry means current directory
				if base == "" {
					base = "."
				}
				cand := filepath.Join(base, target)
				if isDir(cand) {
					target = cand
					cdpathHit = true
					break
				}
			}
		}
	}

	// remember old pwd
	oldPwd, _ := os.Getwd()

	// chdir
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: %v\n", target, err)
		return
	}

	// resolve absolute new pwd
	newPwd, err := os.Getwd()
	if err == nil {
		_ = os.Setenv("OLDPWD", oldPwd)
		_ = os.Setenv("PWD", newPwd)
	}

	// bash prints the path when CDPATH was used and target not under "."
	if cdpathHit {
		fmt.Println(newPwd)
	}
}

// helpers

func getenvHome() string {
	// Prefer HOME; fall back to os.UserHomeDir
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return ""
}

func expandTilde(p string) string {
	if p == "" || p[0] != '~' {
		return p
	}
	// Only expand current user forms: "~" or "~/..."
	if len(p) == 1 || p[1] == '/' || p[1] == '\\' {
		if h := getenvHome(); h != "" {
			return filepath.Join(h, p[1:])
		}
	}
	// "~user" not supported in this minimal shell; return as-is
	return p
}

func splitList(s string) []string {
	// PATH-style splitting; on Unix it's ":", on Windows ";"
	sep := string(os.PathListSeparator)
	return strings.Split(s, sep)
}

func isDir(p string) bool {
	if p == "" {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}
