package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"worker/adb"
	"worker/afk"
	"worker/bot"
	"worker/cfg"

	"github.com/charmbracelet/bubbles/list"
)

func intInput(maxi int) int {
	reader := bufio.NewReader(os.Stdin)
	bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	r := strings.Trim(text, "\r")
	dig, e := strconv.Atoi(r)
	if e != nil || dig > maxi {
		return 0
	}
	return dig
}

func getDevices() []list.Item {
	var devs []list.Item
	d, e := adb.Devices()
	if e != nil {
		return devs
	}
	for _, v := range d {
		descu := fmt.Sprintf("State: %s, T_Id: %v, WMsize: %v", v.DevState, v.TransportId, v.Resolution)
		devs = append(devs, item{title: v.Serial, desc: descu, children: func(m *menuModel) { m.devstatus = runConnect(m) }})
	}
	return devs
}

func runTask(m *menuModel) {
	cf := DtoCfg(m.opts)
	dev, _ := adb.Connect(cf.DeviceSerial)
	gm := afk.New(cf.UserProfile)
	b := bot.New(dev, gm)
	log.Warnf(yellow("CHOSEN RUNTASK >>> %v <<<"), m.choice)
	switch m.choice {
	case "Run all":
		b.UpAll()
	case "Do daily?":
		go func() {
			b.Daily()
		}()
	case "Push Campain?":
		go func() {
			t := b.Task(afk.DOPUSHCAMP)
			b.React(t)
		}()
	case "Kings Tower":
		go func() {
			t := b.Task(afk.Kings)
			b.React(t)
		}()
	case "Towers of Light":
		go func() {
			t := b.Task(afk.Light)
			b.React(t)
		}()
	case "Brutal Citadel":
		go func() {
			t := b.Task(afk.Mauler)
			b.React(t)
		}()
	case "World Tree":
		go func() {
			t := b.Task(afk.Wilder)
			b.React(t)
		}()
	case "Forsaken Necropolis":
		go func() {
			t := b.Task(afk.Graveborn)
			b.React(t)
		}()
	}
}

func Connect(s string) string {
	dev, e := adb.Connect(s)
	if e != nil {
		return ""
	}

	return dev.Serial
}

func runBluestacks(m *menuModel) bool {
	_ = DtoCfg(m.opts)
	e := cfg.RunBlue()
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v", e, "run bluestacks")
		return false
	}
	return true
}

func runConnect(m *menuModel) bool {
	devs := Connect(m.opts[connection])
	if devs != "" {
		return true
	}
	return false
}

func updateDto(v map[string]string) {
	o := DtoCfg(v)
	cfg.Save(cfg.UsrDir(o.UserProfile.Account+".yaml"), o)
}

func CfgDto(conf *cfg.AppConfig) map[string]string {
	dto := make(map[string]string, 0)
	dto[connection] = conf.DeviceSerial
	dto[account] = conf.UserProfile.Account
	dto[game] = conf.UserProfile.Game
	dto[imagick] = strings.Join(conf.Imagick, " ")
	dto[tesseract] = strings.Join(conf.Tesseract, " ")
	dto[bluestacks] = strings.Join(conf.Bluestacks, " ")
	dto[adbp] = reqsoft(adbp, conf)
	dto[magick] = reqsoft(magick, conf)
	dto[bluestacksexe] = reqsoft(bluestacksexe, conf)
	dto[tesserexe] = reqsoft(tesseract, conf)
	// dto[] = conf.
	// dto[] = conf.
	return dto
}

func reqsoft(prop string, cfgn *cfg.AppConfig) string {
	for _, v := range cfgn.RequiredInstalledSoftware {
		if strings.Contains(v, prop) {
			return v
		}
	}
	return ""
}

func DtoCfg(m map[string]string) *cfg.AppConfig {
	res := cfg.Env

	for k, v := range m {
		switch k {
		case connection:
			res.DeviceSerial = v
		case account:
			res.UserProfile.Account = v
		case game:
			res.UserProfile.Game = v
		case imagick:
			res.Imagick = strings.Split(v, " ")
		case tesseract:
			res.Tesseract = strings.Split(v, " ")
		case bluestacks:
			res.Bluestacks = strings.Split(v, " ")
		case adbp, magick, bluestacksexe, tesserexe:
			res.RequiredInstalledSoftware = append(res.RequiredInstalledSoftware, v)
		}
	}
	return res
}
