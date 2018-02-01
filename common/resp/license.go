package resp

// License https://api.github.com/repos/tosone/atom-parsing
type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}
