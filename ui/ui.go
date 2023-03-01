package ui

import (
	"os"
	"worker/cfg"
	"worker/ocr"

	tea "github.com/charmbracelet/bubbletea"
)

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

func init() {
	log = cfg.Logger()
}

func RunMainMenu(c *cfg.Profile) error {
	log.Debug("entered UI")
	options := userSettings(c)
	img := ocrSettings(c, &ocr.Magick{})
	tess := ocrSettings(c, &ocr.Tesseract{})
	m := initialMenuModel(tess, img, options)
	m.list.Title = header
	m.list.SetSize(110, 28)
	m.list.SetShowHelp(true)
	m.list.SetShowPagination(true)
	m.list.SetShowTitle(true)
	m.list.SetShowStatusBar(false)
	m.list.Styles.Title = titleStl
	m.list.Styles.TitleBar = titbarStl

	log.Debugf("Run p, w/ param %s", m)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program:%v", err)
		os.Exit(1)
		return err
	}

	return nil
}
