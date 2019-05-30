package models

import (
	"sort"
)

// import "fmt"

/*
TracksController stores all the tracks information
*/
type TracksController struct {
	Audios    Tracks
	Videos    Tracks
	Subtitles Tracks
}

/*
Tracks is the definition for the slice of *Track
*/
type Tracks []TrackController

/*
BuildTracks creates a new TracksController instance
*/
func BuildTracks(inputs []string) (tracks TracksController) {
	for pos, input := range inputs {
		info, err := GetFileInfo(input)
		if err != nil {
			panic(err)
		}

		info.SetPosition(pos)

		for _, track := range info.Tracks {
			track.SetInfo(&info)
			switch track.Track.Type {
			case "audio":
				tracks.Audios = append(tracks.Audios, track)
			case "video":
				tracks.Videos = append(tracks.Videos, track)
			case "subtitles":
				tracks.Subtitles = append(tracks.Subtitles, track)
			}
		}
	}

	return
}

/*
Filter allows you to filter tracks by a given condition
*/
func (t Tracks) Filter(test func(TrackController) bool) (ret Tracks) {
	for _, item := range t {
		if test(item) {
			ret = append(ret, item)
		}
	}

	return
}

/*
HEVCFilter is used by filter function to filter HEVC sources
*/
func HEVCFilter(track TrackController) bool {
	return track.Track.Codec == "MPEG-H/HEVC/h.265"
}

func videoSorting(videos Tracks) Tracks {
	sort.Slice(videos, func(i, j int) bool {
		// In case they're of the same sice, use position priority setting
		if videos[i].Track.GetHeight() == videos[j].Track.GetHeight() {
			return videos[i].Input.Position > videos[j].Input.Position
		}
		// Take the biggest pixel density image
		return videos[i].Track.GetHeight() < videos[j].Track.GetHeight()
	})

	return videos
}

/*
GetBestVideo returns a pointer to the best available video source track
*/
func (t *TracksController) GetBestVideo() (video *TrackController) {
	videos := t.Videos
	// Try to find-out HEVC sources.
	hevc := videos.Filter(HEVCFilter)

	// If we found just one, return it
	if len(hevc) == 1 {
		return &hevc[0]
	}

	// Don't fuck your brain here, we just wanna sort hevc results in case there
	// are results, otherwise what we want to sort is the original videos input
	if len(hevc) > 0 {
		videos = hevc
	}

	// Sort videos by quality and priority
	videos = videoSorting(videos)

	// Best video should be first
	video = &videos[0]

	return
}

/*
GetBestAudios returns a list with the best available audio source tracks for
the defined languages
*/
func (t *TracksController) GetBestAudios(languages []string) (tracks Tracks) {
	if len(languages) == 0 {
		return t.Audios
	}

	for _, language := range languages {
		resulting := t.GetBestAudio(language)
		tracks = append(tracks, resulting)
	}

	return tracks
}

/*
GetBestAudio among all tracks for the specified language.
*/
func (t *TracksController) GetBestAudio(language string) TrackController {
	audios := t.Audios.Filter(func(track TrackController) bool {
		return track.Track.Properties.Language == language
	})

	if len(audios) == 1 {
		return audios[0]
	}

	if len(audios) == 0 {
		audios = t.Audios
	}

	filtered := t.Audios.Filter(func(track TrackController) bool {
		return track.Track.Properties.CodecID == "A_AAC"
	})

	if len(filtered) == 1 {
		return filtered[0]
	}

	filtered = t.Audios.Filter(func(track TrackController) bool {
		return track.Track.Properties.CodecID == "A_VORBIS"
	})

	if len(filtered) == 1 {
		return filtered[0]
	}

	filtered = t.Audios.Filter(func(track TrackController) bool {
		return track.Track.Properties.CodecID == "A_AC3"
	})

	if len(filtered) == 1 {
		return filtered[0]
	}

	// At the end, if there's no other audio we like, return the one with more priority.
	sort.Slice(audios, func(i, j int) bool {
		return audios[i].Input.Position > audios[j].Input.Position
	})

	return audios[0]
}

/*
GetBestSubtitles among all the tracks, based on given languages and custom definitions.
*/
func (t *TracksController) GetBestSubtitles(languages []string) (subtitles Tracks) {
	if len(languages) == 0 {
		return t.Subtitles
	}

	for _, language := range languages {
		best := t.GetBestSubtitlesForLanguage(language)
		for _, subtitle := range best {
			subtitles = append(subtitles, subtitle)
		}
	}

	return
}

/*
GetBestSubtitlesForLanguage returns all the subtitle tracks for the given language.
Note that subtitles are always return as list, as it may contain forced or
*/
func (t *TracksController) GetBestSubtitlesForLanguage(language string) (subtitles Tracks) {
	subtitles = t.Subtitles.Filter(func(track TrackController) bool {
		return track.Track.Properties.Language == language
	})

	if len(subtitles) == 1 {
		return subtitles
	}

	// Search for forced tracks
	forced := t.Subtitles.Filter(func(track TrackController) bool {
		return track.Track.Properties.Forced
	})

	if len(forced) == 0 {
		// Consider default one is not properly tagged as forced
		def := t.Subtitles.Filter(func(track TrackController) bool {
			return track.Track.Properties.Default
		})

		if len(def) > 0 {
			def[0].Track.Properties.Forced = true
		}
	}

	if len(subtitles) > 2 {
		return subtitles[:1]
	}

	return
}
