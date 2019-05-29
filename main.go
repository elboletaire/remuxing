package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"remuxing/models"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/logrusorgru/aurora"
)

const gray = 13

func merge(
	output string,
	languages []string,
	video models.TrackController, audios models.Tracks, subtitles models.Tracks,
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
		"-d", fmt.Sprint(video.Track.ID),
		// Video route
		video.Input.FileName,
	}

	for i, audio := range audios {
		if i == 0 {
			args = append(args, "--default-track", fmt.Sprint(audio.Track.ID))
		}

		// Ensure audio stream has language set
		args = append(
			args,
			"--language", argIDLabel(audio.Track.ID, audio.Track.Properties.Language),
		)
		// Copy this audio stream
		args = append(args, "-a", fmt.Sprint(audio.Track.ID))
		// Remove its file name
		args = append(args, "--track-name", argIDLabel(audio.Track.ID, ""))
		// Do not copy videos from this file
		args = append(args, "-D")
		// Do not copy subtitles from this file
		args = append(args, "-S")
		// Set the filepath
		args = append(args, audio.Input.FileName)
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

	if len(lang) > 0 {
		languages = strings.Split(lang, ",")
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

func main() {
	_, inputs, languages, _ := parseArgs()
	// output, inputs, languages, verbose := parseArgs()

	// video, audios, subtitles := extractTracks(inputs, languages, verbose)
	tracks := models.BuildTracks(inputs)
	video := tracks.GetBestVideo()
	audios := tracks.GetBestAudios(languages)
	subtitles := tracks.GetBestSubtitles(languages)

	// fmt.Printf("%+v\n", tracks)
	// fmt.Println("Resulting tracks:")
	// pp.Println(tracks)
	fmt.Println("Videos:")
	pp.Println(video.Track)
	fmt.Println("Audios:")
	// pp.Println(audios)
	for _, audio := range audios {
		// pp.Println(audio.Input.FileName)
		pp.Println(audio.Track)
	}
	fmt.Println("Subtitles:")
	// pp.Println(subtitles)
	for _, subtitle := range subtitles {
		pp.Println(subtitle.Track)
	}

	// result := merge(output, languages, video, audios, subtitles)

	// if verbose {
	// 	fmt.Println(string(result))
	// }
}
