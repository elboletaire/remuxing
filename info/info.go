/*
Package info is used for the video information extracted using mkvmerge -i
*/
package info

type properties struct {
	CodecID    string `json:"codec_id"`
	Dimensions string `json:"display_dimensions"`
	Language   string `json:"language"`
}

type track struct {
	ID         int        `json:"id"`
	Type       string     `json:"type"`
	Codec      string     `json:"codec"`
	Properties properties `json:"properties"`
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
	Tracks    []track `json:"tracks"`
}
