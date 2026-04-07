package builtin

import (
	"fmt"
	"os"
	"sort"
	// "strconv"
	"github.com/g3rzi/ctenterd/internal"
)

func init() { Register("ps", PS, "List running processes") }

func PS(args []string) {
	processes, err := internal.GetProcessList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting process list: %v\n", err)
		return
	}
	
	// Sort by PID
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].PID < processes[j].PID
	})
	
	// Print header
	fmt.Printf("%-8s %-8s %-8s %-20s %s\n", "PID", "PPID", "STATE", "NAME", "CMD")
	fmt.Printf("%-8s %-8s %-8s %-20s %s\n", "---", "----", "-----", "----", "---")
	
	for _, proc := range processes {
		fmt.Printf("%-8d %-8d %-8s %-20s %s\n",
			proc.PID,
			proc.PPID,
			proc.State,
			truncate(proc.Name, 20),
			truncate(proc.Cmdline, 40))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
