package piaotong

import (
	"encoding/json"
	"strings"
	"time"
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

func (c *Client) buildRequest(content any) (*Request, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	encrypted, err := c.encrypt(data)
	if err != nil {
		return nil, err
	}

	serialNo, err := c.GenerateSerialNo()
	if err != nil {
		return nil, err
	}

	req := &Request{
		PlatformCode: c.platformCode,
		SignType:     "RSA",
		Format:       "JSON",
		Timestamp:    time.Now().In(timezoneBeijing).Format("2006-01-02 15:04:05"),
		Version:      "1.0",
		SerialNo:     serialNo,
		Content:      encrypted,
	}
	sign, err := c.sign(c.buildSignatureContent(req))
	if err != nil {
		return nil, err
	}

	req.Sign = sign

	return req, nil
}

func (c *Client) buildSignatureContent(req *Request) string {
	pairs := [][2]string{
		{"content", req.Content},
		{"format", req.Format},
		{"platformCode", req.PlatformCode},
		{"serialNo", req.SerialNo},
		{"signType", req.SignType},
		{"timestamp", req.Timestamp},
		{"version", req.Version},
	}

	parts := make([]string, 0, len(pairs))
	for i := range pairs {
		parts = append(parts, strings.Join(pairs[i][:], "="))
	}

	return strings.Join(parts, "&")
}
