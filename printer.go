package main

import (
	"fmt"
	"os"
	"remuxing/models"
	"strings"

	"github.com/logrusorgru/aurora"
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

	fmt.Print(s)
}

func syntaxError(err string) {
	fmt.Println(fmt.Sprintf("syntax error: %s", err))
	os.Exit(1)
}

func title(text string) {
	fmt.Printf(
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
	fmt.Printf(
		aurora.Green("- Track ID %d (%s in %s) from file %s\n").String(),
		track.Track.ID,
		track.Track.Codec,
		track.Track.Properties.Language,
		track.Input.FileName,
	)
}
