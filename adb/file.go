package adb

const (

	remotedir = "/sdcard/"

)

func init() {

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
