/*
Package info is used for the video information extracted using mkvmerge -i
*/
package info

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora"
)

type properties struct {
	CodecID    string `json:"codec_id"`
	Dimensions string `json:"display_dimensions"`
	Language   string `json:"language"`
}

/*
Track for each video, audio or subtitle inside a container.
*/
type Track struct {
	ID         uint       `json:"id"`
	Type       string     `json:"type"`
	Codec      string     `json:"codec"`
	Properties properties `json:"properties"`
	Parent     *Info
}

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
	Tracks    []Track `json:"tracks"`
	FileName  string  `json:"file_name"`
	Position  int
	FileSize  int64
}

/*
SetParent ...
*/
func (track *Track) SetParent(parent *Info) *Info {
	track.Parent = parent

	return track.Parent
}

/*
GetHeight from dimensions
*/
func (track *Track) GetHeight() string {
	return strings.Split(track.Properties.Dimensions, "x")[1]
}

/*
SetPosition for the input queue priority order
*/
func (information *Info) SetPosition(position int) int {
	information.Position = position

	return information.Position
}

/*
SetParents does it for you
*/
func (information *Info) SetParents() *Info {
	for i, track := range information.Tracks {
		track.SetParent(information)
		information.Tracks[i] = track
	}

	return information
}

func filterTracks(items []Track, test func(Track) bool) (ret []Track) {
	for _, item := range items {
		if test(item) {
			ret = append(ret, item)
		}
	}

	return
}

func filterHEVC(track Track) bool {
	return track.Codec == "MPEG-H/HEVC/h.265"
}

/*
GetFileInfo does.. well, that :_)
*/
func GetFileInfo(file string) (information Info) {
	output, mergeErr := exec.Command(
		"mkvmerge",
		"-F",
		"json",
		"-i",
		file,
	).CombinedOutput()

	if mergeErr != nil {
		fmt.Println(fmt.Sprint(mergeErr) + ": " + string(output))
		return
	}

	if err := json.Unmarshal(output, &information); err != nil {
		panic(err)
	}

	fi, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	information.FileSize = fi.Size()

	// We want it in seconds, so we don't care about decimals
	information.Container.Properties.Duration = information.Container.Properties.Duration / 1000 / 1000 / 1000

	return information
}

func videoSorting(videos []Track) []Track {
	sort.Slice(videos, func(i, j int) bool {
		// In case they're of the same sice, use position priority setting
		if videos[i].GetHeight() == videos[j].GetHeight() {
			return videos[i].Parent.Position > videos[j].Parent.Position
		}
		// Take the biggest pixel density image
		return videos[i].GetHeight() > videos[j].GetHeight()
	})

	return videos
}

/*
DecideBestVideo gets the preferred video track from an array of tracks
*/
func DecideBestVideo(videos []Track) Track {
	// Try to find-out HEVC sources.
	hevc := filterTracks(videos, filterHEVC)

	// If we found just one, return it
	if len(hevc) == 1 {
		return hevc[0]
	}

	// Don't fuck your brain here, we just wanna sort hevc results in case there
	// are results, otherwise what we want to sort is the original videos input
	if len(hevc) > 0 {
		videos = hevc
	}

	// Sort videos by quality and priority
	videos = videoSorting(videos)

	// Best video should be first
	return videos[0]
}

/*
DecideBestAudio bamong all tracks for the specified language.
*/
func DecideBestAudio(audios []Track, language string) Track {
	languages := filterTracks(audios, func(track Track) bool {
		return track.Properties.Language == language
	})

	if len(languages) == 0 {
		fmt.Print(aurora.Sprintf(aurora.Red("No audio tracks were found for language %s\n"), language))
		os.Exit(1)
	}

	if len(languages) == 1 {
		return languages[0]
	}

	filtered := filterTracks(languages, func(track Track) bool {
		return track.Properties.CodecID == "A_AAC"
	})

	if len(filtered) == 1 {
		return filtered[0]
	}

	filtered = filterTracks(languages, func(track Track) bool {
		return track.Properties.CodecID == "A_VORBIS"
	})

	if len(filtered) == 1 {
		return filtered[0]
	}

	filtered = filterTracks(languages, func(track Track) bool {
		return track.Properties.CodecID == "A_AC3"
	})

	if len(filtered) == 1 {
		return filtered[0]
	}

	// At the end, if there's no other audio we like, return the one with more priority.
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Parent.Position > languages[j].Parent.Position
	})

	return languages[0]
}

/*
DecideBestAudios among all the tracks using the specified languages.
*/
func DecideBestAudios(audios []Track, languages []string) []Track {
	var tracks []Track
	for _, language := range languages {
		resulting := DecideBestAudio(audios, language)
		tracks = append(tracks, resulting)
	}

	return tracks
}
