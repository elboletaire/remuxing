package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"remuxing/info"
	"strings"

	"github.com/logrusorgru/aurora"
)

func merge() {
	exec.Command(
		"mkvmerge",
		"--title",
		"\"\"",
		"-T",
	)
}

func syntaxError(err string) {
	fmt.Println(fmt.Sprintf("syntax error: %s", err))
	os.Exit(1)
}

func parseArgs() (output string, inputs []string, languages []string) {
	flag.StringVar(&output, "output", "", "The output folder")
	var lang string
	flag.StringVar(&lang, "languages", "", "The desired output languages")

	flag.Parse()

	if len(output) == 0 {
		syntaxError("-output path missing")
	}

	inputs = flag.Args()

	if len(inputs) < 2 {
		syntaxError("at least two inputs are expected")
	}

	languages = strings.Split(lang, ",")

	if len(lang) == 0 || len(languages) < 1 {
		syntaxError("at least one language was expected")
	}

	return
}

func argIDLabel(id uint, lab string) string {
	result := []string{
		fmt.Sprint(id),
		":",
		lab,
	}

	return strings.Join(result, "")
}

func printBestTrack(track info.Track) {
	fmt.Printf(
		aurora.Green("- Track ID %d (%s) from file %s\n").String(),
		track.ID,
		track.Properties.Language,
		track.Parent.FileName,
	)
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

func main() {
	output, inputs, languages := parseArgs()

	var infos []info.Info
	var videos []info.Track
	var audios []info.Track
	var subtitles []info.Track

	title("INPUTS")
	for pos, input := range inputs {
		fmt.Printf(
			aurora.Green("%d. Getting file info for file: %s\n").String(),
			pos+1,
			input,
		)
		information := info.GetFileInfo(input)
		information.SetParents()
		information.SetPosition(pos)

		if !information.Container.Supported {
			fmt.Println(aurora.Sprintf(aurora.Red("File container is not supported for file %s"), input))
			continue
		}

		for _, track := range information.Tracks {
			switch track.Type {
			case "video":
				videos = append(videos, track)
			case "audio":
				audios = append(audios, track)
			case "subtitles":
				subtitles = append(subtitles, track)
			}
		}

		infos = append(infos, information)
	}

	title("VIDEO")
	bestVideo := info.DecideBestVideo(videos)

	printBestTrack(bestVideo)

	title("AUDIO")
	bestAudios := info.DecideBestAudios(audios, languages)
	for _, audio := range bestAudios {
		printBestTrack(audio)
	}

	args := []string{
		"-o", output,
		// Empty name (just in case, should be configurable tho)
		"--title", "",
		// Do not copy audio from the video source
		"-A",
		// Do not copy tracks info from this file
		"-T",
		// Do not copy subtitles either
		"-S",
		// Specify video id to be copied
		"-d", fmt.Sprint(bestVideo.ID),
		// Video route
		bestVideo.Parent.FileName,
	}

	for i, audio := range bestAudios {
		if i == 0 {
			args = append(args, "--default-track", fmt.Sprint(audio.ID))
		}

		// Ensure audio stream has language set
		args = append(
			args,
			"--language", argIDLabel(audio.ID, audio.Properties.Language),
		)
		// Copy this audio stream
		args = append(args, "-a", fmt.Sprint(audio.ID))
		// Remove its file name
		args = append(args, "--track-name", argIDLabel(audio.ID, ""))
		// Do not copy videos from this file
		args = append(args, "-D")
		// Do not copy subtitles from this file
		args = append(args, "-S")
		// Set the filepath
		args = append(args, audio.Parent.FileName)
	}

	title("COMMAND")
	fmt.Printf(aurora.Gray(15, "$ mkvmerge %s\n").String(), strings.Join(args, " "))

	title("RESULT")

	result, err := exec.Command("mkvmerge", args...).CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(result))
		return
	}

	fmt.Println(string(result))

	// fmt.Println(fmt.Sprintf("Output is %s", output))
	// fmt.Println(fmt.Sprintf("Inputs are %s", strings.Join(inputs[:], ",")))
	// fmt.Println(fmt.Sprintf("Languages are %s", strings.Join(languages[:], ",")))
	// args := os.Args[1:]

	// 1. Check expected params are received
	//   - At least two inputs are received
	//   - Output path flag (-o) is also mandatory
}
