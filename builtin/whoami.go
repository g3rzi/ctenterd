package builtin

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	// "syscall"
)

func init() {
	Register("whoami", Whoami, "Show current user (UID/GID)")
}

func Whoami(args []string) {
	uid := os.Getuid()
	gid := os.Getgid()
	
	// Try to get username from /etc/passwd
	username := getUsernameFromPasswd(uid)
	if username == "" {
		username = strconv.Itoa(uid)
	}
	
	groupname := getGroupnameFromGroup(gid)
	if groupname == "" {
		groupname = strconv.Itoa(gid)
	}
	
	fmt.Printf("uid=%d(%s) gid=%d(%s)\n", uid, username, gid, groupname)
}

func getUsernameFromPasswd(uid int) string {
	data, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, ":")
		if len(parts) >= 3 {
			if userUID, err := strconv.Atoi(parts[2]); err == nil && userUID == uid {
				return parts[0]
			}
		}
	}
	
	return ""
}

func getGroupnameFromGroup(gid int) string {
	data, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, ":")
		if len(parts) >= 3 {
			if groupGID, err := strconv.Atoi(parts[2]); err == nil && groupGID == gid {
				return parts[0]
			}
		}
	}
	
	return ""
}
