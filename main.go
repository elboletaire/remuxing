package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/elboletaire/remuxing/models"
)

const gray = 13

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

func main() {
	output, inputs, languages, verbose := parseArgs()

	tracks := models.BuildTracks(inputs)

	video := tracks.GetBestVideo()
	audios := tracks.GetBestAudios(languages)
	subtitles := tracks.GetBestSubtitles(languages)

	command := CommandArguments(output, video, audios, subtitles)

	if verbose {
		title("VIDEOS")
		printTrack(video)
		printTracks("AUDIOS", audios)
		printTracks("SUBTITLES", subtitles)
		printCommand(command)
	}

	result, err := Command(command)

	if err != nil {
		panic(fmt.Sprint(err) + ": " + string(result))
	}

	if verbose {
		title("OUTPUT")
		fmt.Println(string(result))
	}
}
