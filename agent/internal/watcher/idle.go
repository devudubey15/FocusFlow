package watcher

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// IsIdle checks if the system has been idle for more than the specified threshold
func IsIdle(threshold time.Duration) (bool, error) {
	out, err := exec.Command("xprintidle").Output()
	if err != nil {
		return false, err
	}

	idleMsStr := strings.TrimSpace(string(out))
	idleMs, err := strconv.ParseInt(idleMsStr, 10, 64)
	if err != nil {
		return false, err
	}

	return time.Duration(idleMs)*time.Millisecond > threshold, nil
}
