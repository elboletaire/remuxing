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

const gray = 13

func merge(
	output string,
	languages []string,
	video info.Track, audios []info.Track, subtitles []info.Track,
) string {
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
		"-d", fmt.Sprint(video.ID),
		// Video route
		video.Parent.FileName,
	}

	for i, audio := range audios {
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

	s := title("COMMAND")
	s += fmt.Sprintf(aurora.Gray(15, "$ mkvmerge %s\n").String(), strings.Join(args, " "))

	s += title("RESULT")

	result, err := exec.Command("mkvmerge", args...).CombinedOutput()

	if err != nil {
		panic(fmt.Sprint(err) + ": " + string(result))
	}

	s += string(result)

	return s
}

func syntaxError(err string) {
	fmt.Println(fmt.Sprintf("syntax error: %s", err))
	os.Exit(1)
}

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

func parseArgs() (
	output string,
	inputs []string,
	languages []string,
	verbose bool,
) {
	flag.StringVar(&output, "output", "", "The output file.")

	var lang string
	flag.StringVar(&lang, "languages", "", "Languages to be taken from inputs. Order matters, first one will be marked as default track.")

	var help bool
	flag.BoolVar(&help, "h", false, "This help.")
	flag.BoolVar(&verbose, "v", false, "Verbose.")
	flag.Parse()

	if help || len(os.Args) == 1 {
		flag.Usage = func() {
			usage()
			flag.PrintDefaults()
		}

		flag.Usage()
		os.Exit(0)
	}

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

func printBestTrack(verbose bool, track info.Track) {
	if !verbose {
		return
	}

	fmt.Printf(
		aurora.Green("- Track ID %d (%s) from file %s\n").String(),
		track.ID,
		track.Properties.Language,
		track.Parent.FileName,
	)
}

func title(text string) string {
	return fmt.Sprintf(
		aurora.Yellow("\n\n  %s\n  %s %s %s\n  %s\n\n").String(),
		strings.Repeat("#", len(text)+4),
		"#",
		text,
		"#",
		strings.Repeat("#", len(text)+4),
	)
}

func extractTracks(inputs []string, languages []string, verbose bool) (
	video info.Track,
	audios []info.Track,
	subtitles []info.Track,
) {

	var videos []info.Track

	if verbose {
		fmt.Print(title("INPUTS"))
	}

	for pos, input := range inputs {
		if verbose {
			fmt.Printf(
				aurora.Green("%d. Getting file info for file: %s\n").String(),
				pos+1,
				input,
			)
		}
		information := info.GetFileInfo(input)
		information.SetParents()
		information.SetPosition(pos)

		if !information.Container.Supported {
			if verbose {
				fmt.Println(aurora.Red("File container is not supported for file %s").String(), input)
			}
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
	}

	if verbose {
		fmt.Print(title("VIDEO"))
	}
	video = info.DecideBestVideo(videos)
	printBestTrack(verbose, video)

	if verbose {
		fmt.Print(title("AUDIO"))
	}
	audios = info.DecideBestAudios(audios, languages)
	for _, audio := range audios {
		printBestTrack(verbose, audio)
	}

	return
}

func main() {
	output, inputs, languages, verbose := parseArgs()

	video, audios, subtitles := extractTracks(inputs, languages, verbose)

	result := merge(output, languages, video, audios, subtitles)

	if verbose {
		fmt.Println(string(result))
	}
}
