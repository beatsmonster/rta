package process

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

var runPS = func() ([]byte, error) {
	return exec.Command("ps", "-ax", "-o", "pid,ppid,comm").Output()
}

func BuildTree() (children map[int][]int, names map[int]string, err error) {
	out, err := runPS()
	if err != nil {
		return nil, nil, err
	}
	return parseTree(out)
}

func parseTree(out []byte) (map[int][]int, map[int]string, error) {
	children := make(map[int][]int)
	names := make(map[int]string)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Scan() // skip header

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}

		pid, err1 := strconv.Atoi(fields[0])
		ppid, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			continue
		}

		children[ppid] = append(children[ppid], pid)
		names[pid] = fields[2]
	}
	return children, names, nil
}

func HasDescendant(rootPID int, target string, children map[int][]int, names map[int]string) bool {
	queue := []int{rootPID}
	for len(queue) > 0 {
		pid := queue[0]
		queue = queue[1:]
		if names[pid] == target {
			return true
		}
		queue = append(queue, children[pid]...)
	}
	return false
}

func DetectClaude(rootPID int) (bool, error) {
	children, names, err := BuildTree()
	if err != nil {
		return false, err
	}
	return HasDescendant(rootPID, "claude", children, names), nil
}
