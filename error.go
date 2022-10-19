package piaotong

import "fmt"

type Error struct {
	Code     string
	Msg      string
	Sign     string
	SerialNo string
	Content  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s, %s", e.Code, e.Msg)
}
