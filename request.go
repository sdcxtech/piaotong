package piaotong

import (
	"strings"
)

type Request struct {
	PlatformCode string `json:"platformCode"`
	SignType     string `json:"signType"`
	Sign         string `json:"sign"`
	Format       string `json:"format"`
	Timestamp    string `json:"timestamp"`
	Version      string `json:"version"`
	SerialNo     string `json:"serialNo"`
	Content      string `json:"content,omitempty"`
}

func (r *Request) SignatureContent() string {
	pairs := [][2]string{
		{"content", r.Content},
		{"format", r.Format},
		{"platformCode", r.PlatformCode},
		{"serialNo", r.SerialNo},
		{"signType", r.SignType},
		{"timestamp", r.Timestamp},
		{"version", r.Version},
	}

	parts := make([]string, 0, len(pairs))
	for i := range pairs {
		parts = append(parts, strings.Join(pairs[i][:], "="))
	}

	return strings.Join(parts, "&")
}
