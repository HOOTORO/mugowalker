package ui

import c "worker/cfg"

const (
	AccField        = "Account"
	ConnectionField = "Device"
	LoglvlField     = "Loglevel"
)

type AppUser struct {
	Connection    string
	AccountName   string
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

func (au *AppUser) Account() string {
	return au.AccountName
}

func (au *AppUser) DevicePath() string {
	return au.Connection
}

type (
	errMsg error
)

func (au *AppUser) String() string {
	return c.F("Account |> %v\nDevice |> %v\nLogLevel |> %v", au.AccountName, au.Connection, c.Cyan(au.Loglvl))
}
