package builtin

import (
	"fmt"
	"os"
	// "path/filepath"
	"sort"
	"time"
	"github.com/g3rzi/ctenterd/internal"
)

func init() {
	Register("ls", LS, "List directory contents")
}

func LS(args []string) {
	path := "."
	showLong := false
	
	// Parse arguments
	for _, arg := range args {
		if arg == "-l" {
			showLong = true
		} else if !showLong && arg[0] != '-' {
			path = arg
		}
	}
	
	entries, err := internal.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: %v\n", err)
		return
	}
	
	// Sort entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	
	if showLong {
		// Long format
		for _, entry := range entries {
			mode := entry.Mode
			size := entry.Size
			modTime := entry.ModTime
			name := entry.Name
			
			// Format permissions
			perms := formatPermissions(mode)
			
			// Format size
			sizeStr := fmt.Sprintf("%8d", size)
			
			// Format time
			timeStr := modTime.Format("Jan _2 15:04")
			if time.Since(modTime) > 6*30*24*time.Hour { // older than 6 months
				timeStr = modTime.Format("Jan _2  2006")
			}
			
			fmt.Printf("%s %8s %s %s\n", perms, sizeStr, timeStr, name)
		}
	} else {
		// Simple format
		for _, entry := range entries {
			if entry.IsDir {
				fmt.Printf("%s/\n", entry.Name)
			} else {
				fmt.Printf("%s\n", entry.Name)
			}
		}
	}
}

func formatPermissions(mode os.FileMode) string {
	const rwx = "rwxrwxrwx"
	bits := fmt.Sprintf("%09b", mode&os.ModePerm)
	perm := make([]byte, 10)
	
	// File type
	switch mode & os.ModeType {
	case os.ModeDir:
		perm[0] = 'd'
	case os.ModeSymlink:
		perm[0] = 'l'
	case os.ModeNamedPipe:
		perm[0] = 'p'
	case os.ModeSocket:
		perm[0] = 's'
	case os.ModeDevice:
		perm[0] = 'c'
	case os.ModeDevice | os.ModeCharDevice:
		perm[0] = 'b'
	default:
		perm[0] = '-'
	}
	
	// Permissions
	for i, bit := range bits {
		if bit == '1' {
			perm[i+1] = rwx[i]
		} else {
			perm[i+1] = '-'
		}
	}
	
	return string(perm)
}
