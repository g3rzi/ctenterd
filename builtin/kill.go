package builtin

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func init() {
	Register("kill", Kill, "Terminate processes by PID")
}


func Kill(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "kill: missing PID arguments\n")
		fmt.Fprintf(os.Stderr, "Usage: kill [-SIGNAL] PID...\n")
		return
	}

	signal := syscall.SIGTERM // default signal
	var pids []string

	// Parse arguments
	for _, arg := range args {
		if arg[0] == '-' && len(arg) > 1 {
			// Parse signal (simplified - just handle common ones)
			switch arg {
			case "-9", "-KILL":
				signal = syscall.SIGKILL
			case "-15", "-TERM":
				signal = syscall.SIGTERM
			case "-1", "-HUP":
				signal = syscall.SIGHUP
			case "-2", "-INT":
				signal = syscall.SIGINT
			default:
				fmt.Fprintf(os.Stderr, "kill: invalid signal '%s'\n", arg)
				return
			}
		} else {
			pids = append(pids, arg)
		}
	}

	for _, pidStr := range pids {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: invalid PID '%s'\n", pidStr)
			continue
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: cannot find process %d: %v\n", pid, err)
			continue
		}

		if err := process.Signal(signal); err != nil {
			fmt.Fprintf(os.Stderr, "kill: cannot kill process %d: %v\n", pid, err)
		}
	}
}