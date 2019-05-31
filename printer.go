package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/elboletaire/remuxing/models"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

func usage() {
	cmd := os.Args[0]

	s := aurora.Yellow("Usage:\n\n").String()
	s += fmt.Sprintf(
		aurora.Gray(
			gray,
			"   %s -output [file] -languages [langs] [inputs...]\n\n",
		).String(),
		cmd,
	)
	s += aurora.Yellow("Usage example:\n\n").String()
	s += fmt.Sprintf(
		aurora.Gray(
			gray,
			"   %s -v -output output.mkv -languages eng,spa input1.mkv input2.mkv input3.mkv",
		).String(),
		cmd,
	)
	s += strings.Repeat("\n", 3)

	fmt.Fprint(colorable.NewColorableStdout(), s)
}

func printCommand(command []string) {
	title("COMMAND")
	fmt.Fprintf(
		colorable.NewColorableStdout(),
		aurora.Gray(15, "$ mkvmerge %s\n").String(), strings.Join(command, " "),
	)
}

func syntaxError(err string) {
	fmt.Println(fmt.Sprintf("syntax error: %s", err))
	os.Exit(1)
}

func title(text string) {
	fmt.Fprintf(
		colorable.NewColorableStdout(),
		aurora.Yellow("\n\n  %s\n  %s %s %s\n  %s\n\n").String(),
		strings.Repeat("#", len(text)+4),
		"#",
		text,
		"#",
		strings.Repeat("#", len(text)+4),
	)
}

func printTracks(text string, tracks models.Tracks) {
	title(text)
	for _, track := range tracks {
		printTrack(&track)
	}
}

func printTrack(track *models.TrackController) {
	fmt.Fprintf(
		colorable.NewColorableStdout(),
		aurora.Green("- Track ID %d (%s in %s) from file %s\n").String(),
		track.Track.ID,
		track.Track.Codec,
		track.Track.Properties.Language,
		track.Input.FileName,
	)
}
