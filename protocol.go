package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

type Command interface {
	Execute()
}

type SetCommand struct {
	key, value []byte
}
func (c *SetCommand) Execute() {

}

type GetCommand struct {
	key []byte
}
func (c *GetCommand) Execute() {

}

type HelloCommand struct {
	value string
}
func (c  *HelloCommand) Execute() {

}

func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))

	rw := resp.NewWriter(buf)

	for k, v := range m {
		rw.WriteString(k)
		rw.WriteString(":" + v)
	}

	return buf.Bytes()
}
