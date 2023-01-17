package cfg

import (
	"bytes"
	"os"
	"os/exec"
)

func StartProc(sw string, args ...string) (int, error) {
	cmd := exec.Command(sw, args...)
	log.Tracef("start cmd: %v\n", cmd.String())
	e := cmd.Start()
	return cmd.Process.Pid, e
}
func RunProc(sw string, args ...string) *exec.Cmd {
	cmd := exec.Command(sw, args...)
	log.Tracef("run cmd: %v\n", cmd.String())
	cmd.Start()
	return cmd
}

func ProcessInfo(pid int) bool {
	_, e := os.FindProcess(pid)
	if e != nil {
		log.Errorf("Process %v doesn't exist", pid)
		return false
	}

	return true
}

func Tasklist(processname string) (string, error) {
	args := []string{"/fi", f("IMAGENAME eq %v*", processname), "/NH"}
	cmd := exec.Command("tasklist", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	e := cmd.Run()
	if e != nil {
		log.Errorf("Tasklist err: %v", e)
		return "", e
	}

	return buf.String(), nil
}
