package piaotong

import (
	"strings"
)

type Response struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Sign     string `json:"sign"`
	SerialNo string `json:"serialNo"`
	Content  string `json:"content,omitempty"`
}

func (r *Response) SignatureContent() string {
	pairs := [][2]string{
		{"code", r.Code},
		{"content", r.Content},
		{"msg", r.Msg},
		{"serialNo", r.SerialNo},
	}

	parts := make([]string, 0, len(pairs))
	for i := range pairs {
		parts = append(parts, strings.Join(pairs[i][:], "="))
	}

	return strings.Join(parts, "&")
}

func (r *Response) isError() bool {
	return r.Code != "0000"
}

func (r *Response) toError() *Error {
	return &Error{
		Code:     r.Code,
		Msg:      r.Msg,
		Sign:     r.Sign,
		SerialNo: r.SerialNo,
		Content:  r.Content,
	}
}
