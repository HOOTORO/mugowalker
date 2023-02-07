package ui

import (
	"os"
	"worker/cfg"
	"worker/ocr"

	"github.com/sirupsen/logrus"

	tea "github.com/charmbracelet/bubbletea"
)

// var red, green, cyan, yellow, mag func(...interface{}) string

// func init() {
// 	red = color.New(color.FgHiRed).SprintFunc()
// 	green = color.New(color.FgHiGreen).SprintFunc()
// 	cyan = color.New(color.FgHiCyan).SprintFunc()
// 	yellow = color.New(color.FgHiYellow).SprintFunc()
// 	mag = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()
// }

type AppUser struct {
	Connection    string
	Account       string
	AndroidGameID string
	VMName        string
	Loglvl        string
}

func (au *AppUser) Loglevel() string {
	return au.Loglvl
}

func (au *AppUser) Game() string {
	return au.AndroidGameID
}

func (au *AppUser) Acccount() string {
	return au.Account
}

func (au *AppUser) DevicePath() string {
	return au.Connection
}

type OcrSettings struct {
}

// keymapping
const (
	connection   = "connection"
	account      = "account"
	game         = "game"
	taskconfigs  = "taskconfigs"
	imagick      = "imagick"
	tesseract    = "tesseract"
	blueInstance = "bluestance"
	bluePackage  = "bluepackage"
	bluexe       = "HD-Player"
)

var log *logrus.Logger

func init() {
	log = cfg.Logger()
}

func RunMainMenu(c *cfg.Profile) error {
	log.Debug("entered UI")
	options := userSettings(c)
	img := ocrSettings(c, ocr.Magick{})
	tess := ocrSettings(c, ocr.Tesseract{})
	m := InitialMenuModel(tess, img, options)
	m.menulist.Title = header
	m.menulist.SetSize(110, 28)
	m.menulist.SetShowHelp(true)
	m.menulist.SetShowPagination(true)
	m.menulist.SetShowTitle(true)
	m.menulist.SetShowStatusBar(false)
	m.menulist.Styles.Title = titleStl
	m.menulist.Styles.TitleBar = titbarStl

	log.Debugf("Run p, w/ param %s", m)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program:%v", err)
		os.Exit(1)
		return err
	}

	return nil
}

func NotifyUI(task, desc string) {

}
