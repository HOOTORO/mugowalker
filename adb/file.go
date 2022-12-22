package adb

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	workdir   = ".adb"
	remotedir = "/sdcard/"
	localdir  = ".afk_data"
)

func init() {
	usr := os.Getenv("USERPROFILE")
	wd := filepath.Join(usr, workdir)
	_, e := os.Lstat(wd)
	if e == nil || os.IsNotExist(e) {
		//		os.MkdirAll(wd, os.ModeDir)
		wd, _ = os.Getwd()
		e := os.MkdirAll(filepath.Join(wd, localdir), os.ModeDir)
		e = os.Chdir(localdir)
		fmt.Printf("\ninit: success; pwd: %v\n\n", wd)
        _, e = os.Lstat("app.log")
        if os.IsNotExist(e){
            _, e = os.Create("app.log")
        }

	} else {
		pwd, _ := os.Getwd()
		fmt.Printf("init: fail; error: %v\npwd will be used: %v\n\n", e.Error(), pwd)
	}
}

// Push Pushes the local file to the remote one.
func (d *Device) Push(local, remote string) error {
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"push", local, remote,
	}}
	_, err := cmd.Call()
	return err
}

// Pull Pulls the remote file to the local one.
func (d *Device) Pull(remote, local string) error {
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"pull", remotedir + remote, local,
	}}
	_, err := cmd.Call()
	return err
}
