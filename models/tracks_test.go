package models

import (
	"testing"

	"github.com/elboletaire/remuxing/tests"
)

func TestGetBestVideoReturnsUniqueHEVCSource(t *testing.T) {
	tracks := TracksController{
		Videos: Tracks{
			TrackController{
				Track: &Track{
					ID:    0,
					Codec: "Whatever other codec",
				},
			},
			TrackController{
				Track: &Track{
					ID:    1,
					Codec: "MPEG-H/HEVC/h.265",
				},
			},
		},
	}

	tests.Equals(t, "1", tracks.GetBestVideo().Track.GetID())
}

func TestGetBestVideoDecidesBetweenHEVCSourcesBasedOnDimensions(t *testing.T) {
	big := "1920x1080"
	mid := "1280x720"
	small := "640x480"

	tracks := TracksController{
		Videos: Tracks{
			TrackController{
				Track: &Track{
					ID:    0,
					Codec: "MPEG-H/HEVC/h.265",
					Properties: properties{
						Dimensions: &big,
					},
				},
			},
			TrackController{
				Track: &Track{
					ID:    1,
					Codec: "MPEG-H/HEVC/h.265",
					Properties: properties{
						Dimensions: &mid,
					},
				},
			},
			TrackController{
				Track: &Track{
					ID:    2,
					Codec: "MPEG-H/HEVC/h.265",
					Properties: properties{
						Dimensions: &small,
					},
				},
			},
		},
	}

	tests.Equals(t, "0", tracks.GetBestVideo().Track.GetID())
}
func TestGetBestVideoDecidesBetweenHEVCSourcesBasedOnPosition(t *testing.T) {
	dimensions := "1920x1080"

	tracks := TracksController{
		Videos: Tracks{
			TrackController{
				Input: &Info{
					Position: 2,
				},
				Track: &Track{
					ID:    0,
					Codec: "MPEG-H/HEVC/h.265",
					Properties: properties{
						Dimensions: &dimensions,
					},
				},
			},
			TrackController{
				Input: &Info{
					Position: 1,
				},
				Track: &Track{
					ID:    1,
					Codec: "MPEG-H/HEVC/h.265",
					Properties: properties{
						Dimensions: &dimensions,
					},
				},
			},
			TrackController{
				Input: &Info{
					Position: 0,
				},
				Track: &Track{
					ID:    2,
					Codec: "MPEG-H/HEVC/h.265",
					Properties: properties{
						Dimensions: &dimensions,
					},
				},
			},
		},
	}

	tests.Equals(t, "0", tracks.GetBestVideo().Track.GetID())
}

func TestGetBestVideoDecidesBetweenNonHEVCSourcesBasedOnPosition(t *testing.T) {
	dimensions := "1920x1080"

	tracks := TracksController{
		Videos: Tracks{
			TrackController{
				Input: &Info{
					Position: 2,
				},
				Track: &Track{
					ID:    0,
					Codec: "whatever",
					Properties: properties{
						Dimensions: &dimensions,
					},
				},
			},
			TrackController{
				Input: &Info{
					Position: 1,
				},
				Track: &Track{
					ID:    1,
					Codec: "whatever",
					Properties: properties{
						Dimensions: &dimensions,
					},
				},
			},
			TrackController{
				Input: &Info{
					Position: 0,
				},
				Track: &Track{
					ID:    2,
					Codec: "whatever",
					Properties: properties{
						Dimensions: &dimensions,
					},
				},
			},
		},
	}

	tests.Equals(t, "0", tracks.GetBestVideo().Track.GetID())
}

func TestGetBestAudioReturnsTheOnlyOneWithSpecifiedLanguage(t *testing.T) {
	tracks := TracksController{
		Audios: Tracks{
			TrackController{
				Track: &Track{
					ID: 0,
					Properties: properties{
						Language: "eng",
					},
				},
			},
			TrackController{
				Track: &Track{
					ID: 1,
					Properties: properties{
						Language: "spa",
					},
				},
			},
		},
	}

	tests.Equals(t, "1", tracks.GetBestAudio("spa").Track.GetID())
}

func TestGetBestAudioDecidesBetweenCodecsWhenMultipleSourcesOfSameLanguage(t *testing.T) {
	tracks := TracksController{
		Audios: Tracks{
			TrackController{
				Track: &Track{
					ID: 0,
					Properties: properties{
						Language: "eng",
						CodecID:  "A_VORBIS",
					},
				},
			},
			TrackController{
				Track: &Track{
					ID: 1,
					Properties: properties{
						Language: "eng",
						CodecID:  "A_AAC",
					},
				},
			},
		},
	}

	tests.Equals(t, "1", tracks.GetBestAudios([]string{"eng"})[0].Track.GetID())

	tracks = TracksController{
		Audios: Tracks{
			TrackController{
				Track: &Track{
					ID: 0,
					Properties: properties{
						Language: "eng",
						CodecID:  "A_VORBIS",
					},
				},
			},
			TrackController{
				Track: &Track{
					ID: 1,
					Properties: properties{
						Language: "eng",
						CodecID:  "A_AC3",
					},
				},
			},
		},
	}

	tests.Equals(t, "0", tracks.GetBestAudios([]string{"eng"})[0].Track.GetID())
}

func TestGetBestAudioDecidesBetweenCodecsBasedOnPositionIfNoKnownCodecsFound(t *testing.T) {
	tracks := TracksController{
		Audios: Tracks{
			TrackController{
				Input: &Info{
					Position: 1,
				},
				Track: &Track{
					ID: 0,
					Properties: properties{
						Language: "eng",
						CodecID:  "whatever",
					},
				},
			},
			TrackController{
				Input: &Info{
					Position: 0,
				},
				Track: &Track{
					ID: 1,
					Properties: properties{
						Language: "eng",
						CodecID:  "whatever",
					},
				},
			},
		},
	}

	tests.Equals(t, "0", tracks.GetBestAudios([]string{"eng"})[0].Track.GetID())
}
