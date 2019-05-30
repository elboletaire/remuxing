/*
Package models is used for the video information extracted using mkvmerge -i
*/
package models

import (
	"encoding/json"
	"os"
	"os/exec"
)

type props struct {
	Duration uint64
}

type container struct {
	Properties props
	Supported  bool
}

/*
Info is the main video information object/struct
*/
type Info struct {
	Container container
	Tracks    Tracks `json:"tracks"`
	FileName  string `json:"file_name"`
	Position  int
	FileSize  int64
}

/*
SetPosition for the input queue priority order
*/
func (information *Info) SetPosition(position int) int {
	information.Position = position

	return information.Position
}

/*
GetFileInfo does.. well, that :_)
*/
func GetFileInfo(file string) (information Info, err error) {
	output, err := exec.Command(
		"mkvmerge",
		"-F",
		"json",
		"-i",
		file,
	).CombinedOutput()

	if err != nil {
		return
	}

	if err = json.Unmarshal(output, &information); err != nil {
		return
	}

	fi, err := os.Stat(file)
	if err != nil {
		return
	}

	information.FileSize = fi.Size()

	// We want it in seconds, so we don't care about decimals
	information.Container.Properties.Duration = information.Container.Properties.Duration / 1000 / 1000 / 1000

	return information, nil
}
