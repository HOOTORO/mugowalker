package emulator

import (
	"errors"
	"fmt"
	"os"
	"worker/cfg"

	"github.com/fatih/color"
)

type Vendor int

const (
	BlueStacks Vendor = iota + 1
)

// args
const (
	instance = "--instance"
	command  = "--cmd"
	pack     = "--package"
)

var (
	f             = fmt.Sprintf
	blue          = color.New(color.FgHiBlue).SprintFunc()
	defins        = "Pi"
	defcmd        = "launchApp"
	afktestclient = "com.lilithgames.hgame.gp.id"
	afkglobal     = "com.lilithgames.hgame.gp"
)

// Emulator runs locally
type Emulator struct {
	pid    int
	state  *os.ProcessState
	online bool
}

// New emulator run, optional args:
// 1. <vmname> to run (example: Rvc64 - it is not name setted in Multi-instance manager, can be find in BlueStacks_nxt\Engine foldername of dir with target vm)
// 2. <package> to start (example: "com.lilithgames.hgame.gp" - afk arena)
func New(em Vendor, fn func(string, string), args ...string) (*Emulator, error) {
	if em != BlueStacks {
		return nil, errors.New("only bluestacks supported for now")
	}
	res := &Emulator{}
	s, err := cfg.Tasklist(cfg.BluestacksExe)
	if err != nil && s == "" {
		go func() {
			cmd := cfg.RunProc(cfg.BluestacksExe)
			cmd.Start()
			res.pid = cmd.Process.Pid
			res.state = cmd.ProcessState
			res.online = true
			fn(blue("BlueStacks"), f("Started! Pid: %v", res.pid))
			e := cmd.Wait()
			if e != nil {
				fn(f("BlueStacks"), f("stopped: %v", e))
			}
			fn(f("BlueStacks"), f("finished: %v", e))
			res.pid = 0
			res.online = false
		}()
	} else {
		res.pid = cfg.ToInt(s)
		res.online = cfg.ProcessInfo(res.pid)
	}
	return res, nil

}

// IsOnline vm?
func (e *Emulator) IsOnline() bool {
	return e.online
}

// Kill using os.Process.Kill under the hood цith all that it implies
func (e *Emulator) Kill() bool {
	return cfg.Kill(e.pid)
}

// func Kill() {
// 	cmd := exec.Command(kill, killargs...)
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	r := cmd.Run()
// 	if r != nil {
// 		panic("emu err")
// 	}
// 	color.HiRed(out.String())
// }
// func Start() {
// 	cmd := exec.Command(vm, vmargs...)
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	r := cmd.Run()
// 	if r != nil {
// 		panic("emu err")
// 	}
// 	color.HiRed(out.String())
// }
