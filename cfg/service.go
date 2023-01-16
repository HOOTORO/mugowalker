package cfg

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func StartProc(sw string, args ...string) (int, error) {
	cmd := exec.Command(sw, args...)
	log.Tracef("run cmd: %v\n", cmd.String())
	e := cmd.Start()
	return cmd.Process.Pid, e
}

func ProcessInfo(pid int) bool {
	_, e := os.FindProcess(pid)
	if e != nil {
		log.Errorf("Process %v doesn't exist", pid)
		return false
	}

	return true
}

func PidInfo(p string) string {
	// args := []string{"ps", "|", "where", "Name", "-Like", f("'%v'", p)}
	args := []string{"/fi", f("IMAGENAME eq %v*", p), "/NH", "/V"}
	cmd := exec.Command("tasklist", args...)
	log.Warnf("cmd: %v", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	r := cmd.Run()
	if r != nil {
		log.Errorf("err:%v", r)
	}
	res := strings.Fields(out.String())
	return strings.Join([]string{res[0], res[1], res[6]}, "#")
}
