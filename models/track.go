package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type properties struct {
	CodecID    string  `json:"codec_id"`
	Dimensions *string `json:"display_dimensions"`
	Language   string  `json:"language"`
	Forced     bool    `json:"forced_track"`
	Default    bool    `json:"default_track"`
}

/*
Track data for each video, audio or subtitle inside a container.
*/
type Track struct {
	ID         uint       `json:"id"`
	Type       string     `json:"type"`
	Codec      string     `json:"codec"`
	Properties properties `json:"properties"`
}

/*
TrackController for each video, audio or subtitle inside a container.
*/
type TrackController struct {
	Input *Info
	Track *Track
}

/*
UnmarshalJSON ensures info is properly extracted to TrackController
*/
func (track *TrackController) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &track.Track)
}

/*
SetInfo ...
*/
func (track *TrackController) SetInfo(input *Info) *Info {
	track.Input = input

	return track.Input
}

/*
GetHeight from dimensions
*/
func (track *Track) GetHeight() string {
	return strings.Split(*track.Properties.Dimensions, "x")[1]
}

/*
GetID returns the track id as string
*/
func (track *Track) GetID() string {
	return fmt.Sprint(track.ID)
}

/*
GetArgIDLabel returns an id argument with the specified label (mkvmerge syntax)
*/
func (track *Track) GetArgIDLabel(label string) string {
	return track.GetArgID() + label
}

/*
GetArgID returns the id argument (mkvmerge syntax)
*/
func (track *Track) GetArgID() string {
	return fmt.Sprint(track.ID) + ":"
}
