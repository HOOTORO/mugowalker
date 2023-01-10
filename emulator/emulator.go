package emulator

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var (
	psarg    = []string{"/fi", "IMAGENAME eq HD-Player.exe"}
	vmargs   = []string{"--instance", "Rvc64", "--cmd", "launchApp", "--package", "com.lilithgames.hgame.gp.id"}
	killargs = []string{"/fi", "/IM", "HD-Player.exe", "/F"}
)

const (
	vm     = "HD-Player"
	kill   = "taskkill"
	psinfo = "TASKLIST"
)

type Emulator struct {
	online bool
}

func IsOnline() bool {
	cmd := exec.Command(psinfo, psarg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	r := cmd.Run()
	if r != nil {
		panic("emu err")
	}
	return strings.Contains(out.String(), vm)
}

func Kill() {
	cmd := exec.Command(kill, killargs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	r := cmd.Run()
	if r != nil {
		panic("emu err")
	}
	color.HiRed(out.String())
}

func Start() {
	cmd := exec.Command(vm, vmargs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	r := cmd.Run()
	if r != nil {
		panic("emu err")
	}
	color.HiRed(out.String())
}
