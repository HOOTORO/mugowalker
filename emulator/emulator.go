package emulator

import (
	"errors"
	"fmt"
	"os"
	"worker/cfg"

	"github.com/fatih/color"
)

type Vendor uint

var vendors = []string{"HD-Player", "Nox"}

const (
	BlueStacks Vendor = iota + 1
	Nox
	AllVendors = BlueStacks | Nox
)

func (v Vendor) String() string {
	return vendors[v-1]
}
func (v Vendor) Values() []string {
	return vendors
}

// args
const (
	instance = "--instance"
	command  = "--cmd"
	pack     = "--package"
)

var (
	f             = fmt.Sprintf
	log           = cfg.Logger()
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
	if em != BlueStacks || em != Nox {
		return nil, errors.New("only bluestacks supported for now")
	}
	res := &Emulator{}
	s, err := cfg.Tasklist(em.String())
	if err != nil {
		go func() {
			cmd := cfg.RunProc(em.String())
			cmd.Start()
			res.pid = cmd.Process.Pid
			res.state = cmd.ProcessState
			res.online = true
			fn(blue(em.String()), f("Started! Pid: %v", res.pid))
			e := cmd.Wait()
			if e != nil {
				fn(f(em.String()), f("stopped: %v", e))
			}
			fn(f(em.String()), f("finished: %v", e))
			res.pid = 0
			res.online = false
		}()
	} else {
		res.pid = s[0].Pid
		res.online = cfg.IsProcess(res.pid)
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
