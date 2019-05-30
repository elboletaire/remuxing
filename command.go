package main

import (
	"fmt"
	"os/exec"
	"github.com/elboletaire/remuxing/models"
)

/*
CommandArguments generates the command line for mkvmerge based on the resulting video, audios and subs.
*/
func CommandArguments(
	output string,
	video *models.TrackController,
	audios models.Tracks,
	subtitles models.Tracks,
) (command []string) {
	// The output line "-o {.filename}"
	command = []string{"-o", output}
	// Video options
	command = videoString(video, command)
	// Audio options
	command = audiosString(audios, command)
	// Subtitles options
	command = subtitlesString(subtitles, command)

	return command
}

/*
Command executes the mkvmerge system command with the given args
*/
func Command(args []string) (result []byte, err error) {
	result, err = exec.Command("mkvmerge", args...).CombinedOutput()

	return
}

func videoString(video *models.TrackController, command []string) []string {
	command = append(
		command,
		// Empty name (just in case, should be configurable tho)
		"--title", "",
		// Do not copy audio from the video source
		"-A",
		// Do not copy tracks info from this file
		"-T",
		// Do not copy subtitles either
		"-S",
		// Specify video id to be copied
		"-d", video.Track.GetID(),
		// Video route
		video.Input.FileName,
	)

	return command
}

func audiosString(audios models.Tracks, command []string) []string {
	for i, audio := range audios {
		command = append(command, "-T")
		// Hardcode first as default (they should come already sorted by priority)
		if i == 0 {
			command = append(command, "--default-track", audio.Track.GetID())
		}

		// Ensure audio stream has language set
		command = append(
			command,
			"--language", audio.Track.GetArgIDLabel(audio.Track.Properties.Language),
			// Copy this audio stream
			"-a", audio.Track.GetID(),
			// Remove its file name
			"--track-name", audio.Track.GetArgID(),
			// Do not copy videos from this file
			"-D",
			// Do not copy subtitles from this file
			"-S",
			// Set the filepath
			audio.Input.FileName,
		)
	}

	return command
}

func subtitlesString(subtitles models.Tracks, command []string) []string {
	for _, subtitle := range subtitles {
		command = append(
			command,
			"-T",
			// Copy this subtitle track
			"-s", fmt.Sprint(subtitle.Track.ID),
			// Remove its file name
			"--track-name", subtitle.Track.GetArgID(),
			// Do not copy audios nor videos from this track
			"-D", "-A",
		)

		if subtitle.Track.Properties.Forced {
			command = append(
				command,
				"--forced-track", subtitle.Track.GetArgIDLabel("true"),
			)
		}

		// The subtitle file source
		command = append(command, subtitle.Input.FileName)
	}

	return command
}
