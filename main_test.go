package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMainRunsStatus(t *testing.T) {
	if os.Getenv("RTA_TEST_MAIN") == "1" {
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainRunsStatus")
	cmd.Env = append(os.Environ(), "RTA_TEST_MAIN=1")
	_ = cmd.Run()
}
