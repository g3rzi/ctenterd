package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	// "path/filepath"
	"strconv"
	"strings"
)

type Process struct {
	PID     int
	PPID    int
	Name    string
	State   string
	Cmdline string
}

func GetProcessList() ([]*Process, error) {
	procDir, err := os.Open("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc: %v", err)
	}
	defer procDir.Close()
	
	entries, err := procDir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc: %v", err)
	}
	
	var processes []*Process
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		pidStr := entry.Name()
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue // not a PID directory
		}
		
		proc, err := parseProcess(pid)
		if err != nil {
			continue // process may have exited
		}
		
		processes = append(processes, proc)
	}
	
	return processes, nil
}

func parseProcess(pid int) (*Process, error) {
	// Read /proc/PID/stat
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	statData, err := ioutil.ReadFile(statPath)
	if err != nil {
		return nil, err
	}
	
	// Parse stat file - format is complex due to comm field potentially containing spaces and parentheses
	statStr := string(statData)
	
	// Find the last ')' to split comm from other fields
	lastParen := strings.LastIndex(statStr, ")")
	if lastParen == -1 {
		return nil, fmt.Errorf("invalid stat format")
	}
	
	// Extract comm (process name)
	firstParen := strings.Index(statStr, "(")
	if firstParen == -1 {
		return nil, fmt.Errorf("invalid stat format")
	}
	comm := statStr[firstParen+1 : lastParen]
	
	// Parse remaining fields
	remainingFields := strings.Fields(statStr[lastParen+1:])
	if len(remainingFields) < 2 {
		return nil, fmt.Errorf("insufficient stat fields")
	}
	
	state := remainingFields[0]
	ppidStr := remainingFields[1]
	
	ppid, err := strconv.Atoi(ppidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ppid: %v", err)
	}
	
	// Read cmdline
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdlineData, err := ioutil.ReadFile(cmdlinePath)
	cmdline := ""
	if err == nil {
		// cmdline is null-separated, convert to space-separated
		cmdline = strings.Replace(string(cmdlineData), "\x00", " ", -1)
		cmdline = strings.TrimSpace(cmdline)
	}
	
	if cmdline == "" {
		cmdline = fmt.Sprintf("[%s]", comm)
	}
	
	return &Process{
		PID:     pid,
		PPID:    ppid,
		Name:    comm,
		State:   state,
		Cmdline: cmdline,
	}, nil
}