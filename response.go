package piaotong

import "encoding/json"

type Response struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Sign     string `json:"sign"`
	SerialNo string `json:"serialNo"`
	Content  string `json:"content,omitempty"`
}

func (c *Client) decryptAndUnmarshalResponse(content string, v any) error {
	data, err := c.Decrypt(content)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
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
