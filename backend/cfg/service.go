package cfg

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	Name string
	Pid  int
}

func StartProc(x string, args ...string) (int, error) {
	cmd := exec.Command(x, args...)
	log.Tracef("start cmd: %v\n", cmd.String())
	e := cmd.Start()
	return cmd.Process.Pid, e
}
func RunProc(x string, args ...string) *exec.Cmd {
	cmd := exec.Command(x, args...)
	log.Tracef("run cmd: %v\n", cmd.String())
	cmd.Start()
	return cmd
}

func IsProcess(pid int) bool {
	_, e := os.FindProcess(pid)
	if e != nil {
		log.Errorf("Process %v doesn't exist", pid)
		return false
	}

	return true
}

func Tasklist(processname string) ([]ProcessInfo, error) {
	args := []string{"/fi", F("IMAGENAME eq %v*", processname), "/NH"}
	cmd := exec.Command("tasklist", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	e := cmd.Run()
	if e != nil {
		log.Errorf("Tasklist err: %v", e)
		return nil, e
	}
	parseTasklist(buf.String())
	return parseTasklist(buf.String()), nil
}

func Kill(pid int) bool {
	prc, err := os.FindProcess(pid)
	if err != nil {
		log.Info("Process alreary killed: ", pid)
		return true
	}
	err = prc.Kill()
	if err != nil {
		log.Errorf("Cannot kill, err: %v", err)
		return false
	}
	return true
}

func Intersect(main []string, compared []string) (result []string) {
	for _, v := range main {
		for _, kw := range compared {
			if strings.Contains(v, kw) {
				result = append(result, v)
			}
		}
	}
	return

}

func Regex(s, r string) (res []uint) {
	re := regexp.MustCompile(r)
	for _, v := range re.FindStringSubmatch(s) {
		i, err := strconv.ParseInt(v, 10, 32)
		if err == nil {
			res = append(res, uint(i))
		}
	}
	return
}

func parseTasklist(s string) []ProcessInfo {

	var res []ProcessInfo
	fields := strings.Fields(s)
	if len(fields) > 0 && fields[0] == "INFO:" {
		return nil
	}
	for i, f := range fields {
		if strings.HasSuffix(f, ".exe") {
			res = append(res, ProcessInfo{Name: f[:len(f)-4], Pid: ToInt(fields[i+1])})
		}
	}
	log.Debugf("PARSETASKS RES: %v |> %v", len(res), res)
	return res

}
