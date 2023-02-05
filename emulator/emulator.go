package emulator

import (
	"errors"
	"os"
	c "worker/cfg"
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
	log           = c.Logger()
	defcmd        = "launchApp"
	afktestclient = "com.lilithgames.hgame.gp.id"
	afkglobal     = "com.lilithgames.hgame.gp"
)

// Emulator runs locally
type Emulator struct {
	v      Vendor
	pid    int
	state  *os.ProcessState
	online bool
}

var runninVMs []*Emulator

// Run emulator, optional args:
// 1. <vmname> to run (example: Rvc64 - it is not name setted in Multi-instance manager, can be find in BlueStacks_nxt\Engine foldername of dir with target vm)
// 2. <package> to start (example: "com.lilithgames.hgame.gp" - afk arena)
func Run(em Vendor, args ...string) (int, error) {
	if em != BlueStacks || em != Nox {
		return 0, errors.New("only bluestacks supported for now")
	}
	res := &Emulator{v: em}
	s, err := c.Tasklist(em.String())
	if err != nil {
		go func() {
			cmd := c.RunProc(em.String())
			cmd.Start()
			res.pid = cmd.Process.Pid
			res.state = cmd.ProcessState
			res.online = true
			log.Trace(c.Blue(em.String()), c.F("Started! Pid: %v", res.pid))
			e := cmd.Wait()
			if e != nil {
				log.Trace(c.F(em.String()), c.F("stopped: %v", e))
			}
			log.Trace(c.F(em.String()), c.F("finished: %v", e))
			res.pid = 0
			res.online = false
		}()
	} else {
		res.pid = s[0].Pid
		res.online = c.IsProcess(res.pid)
	}
	return res.pid, nil

}

// IsOnline vm?
func (e *Emulator) IsOnline() bool {
	return e.online
}

// Kill using os.Process.Kill under the hood цith all that it implies
func (e *Emulator) Kill() bool {
	return c.Kill(e.pid)
}

func Emu(v Vendor) *Emulator {
	for _, em := range runninVMs {
		if em.v == v {
			return em
		}
	}
	return nil
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
