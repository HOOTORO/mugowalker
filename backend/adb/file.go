package adb

import (
	"errors"
	"fmt"
	"mugowalker/backend/localstore"
	"path"
	"strconv"
)

const (
	remotedir = "/sdcard/"
)

type RemoteFile struct {
	name  string
	path  string
	local string
}

// Push Pushes the local file to the remote one.
func (d *Device) Push(local, remote string) error {
	cmd := Cmd{Args: []string{
		"-t", strconv.Itoa(d.TransportId),
		"push", local, remote,
	}}
	_, err := cmd.Call()
	return err
}

// Pull Pulls the remote file to the local one.
func (d *Device) Pull() (string, error) {
	if len(d.Files) > 0 {
		lastFile := d.Files[len(d.Files)-1]
		cmd := Cmd{Args: []string{
			"-t", strconv.Itoa(d.TransportId),
			"pull", lastFile.path, localstore.TempDir(),
		}}
		result, err := cmd.Call()
		if err != nil {
			return fmt.Sprintf("PULL call: %v, ERROR: %v", result, err), err
		}
		lastFile.local = path.Join(localstore.TempDir(), lastFile.name)
		return lastFile.local, nil
	} else {
		return "PULL ERROR:", errors.New("Remote device haven't any files")
	}
}
