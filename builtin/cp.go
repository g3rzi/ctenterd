package builtin

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func init() {
	Register("cp", Copy, "Copy files or directories")
}


// cp - Copy files or directories
func Copy(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "cp: missing file arguments\n")
		fmt.Fprintf(os.Stderr, "Usage: cp SOURCE DEST\n")
		return
	}

	recursive := false
	var sources []string
	var dest string

	// Parse arguments
	for i, arg := range args {
		if arg == "-r" || arg == "-R" {
			recursive = true
		} else if i == len(args)-1 {
			dest = arg
		} else if arg[0] != '-' {
			sources = append(sources, arg)
		}
	}

	if len(sources) == 0 {
		fmt.Fprintf(os.Stderr, "cp: no source files specified\n")
		return
	}

	// Handle multiple sources
	if len(sources) > 1 {
		// Destination must be a directory
		destInfo, err := os.Stat(dest)
		if err != nil || !destInfo.IsDir() {
			fmt.Fprintf(os.Stderr, "cp: target '%s' is not a directory\n", dest)
			return
		}

		for _, src := range sources {
			basename := filepath.Base(src)
			destPath := filepath.Join(dest, basename)
			if err := copyPath(src, destPath, recursive); err != nil {
				fmt.Fprintf(os.Stderr, "cp: cannot copy '%s' to '%s': %v\n", src, destPath, err)
			}
		}
	} else {
		// Single source
		if err := copyPath(sources[0], dest, recursive); err != nil {
			fmt.Fprintf(os.Stderr, "cp: cannot copy '%s' to '%s': %v\n", sources[0], dest, err)
		}
	}
}

func copyPath(src, dest string, recursive bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if !recursive {
			return fmt.Errorf("omitting directory '%s' (use -r to copy directories)", src)
		}
		return copyDir(src, dest)
	}

	return copyFile(src, dest)
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	return os.Chmod(dest, srcInfo.Mode())
}

func copyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}