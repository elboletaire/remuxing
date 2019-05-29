package models

import (
	"encoding/json"
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
	type Temporal Track

	if err := json.Unmarshal(data, &track.Track); err != nil {
		return err
	}

	return nil
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
