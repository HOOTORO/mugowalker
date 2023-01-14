package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"worker/adb"
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
		devs = append(devs, menuItem("No devices found, try to connect"))
		return devs
	}
	for _, v := range d {
		devs = append(devs, menuItem(v.Serial))
	}
	return devs
}

func Connect(d string) []list.Item {
	var items []list.Item
	dev, e := adb.Connect(d)
	if e != nil {
		items = append(items, menuItem("Connection failed"))
	}
	items = append(items, menuItem(dev.Serial))
	return items
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

func runBluestacks(m menuModel) {
	_ = DtoCfg(m.opts)
	e := cfg.RunBlue()
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v", e, "run bluestacks")
	}
}
