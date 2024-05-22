package client

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)


type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

// 模拟客户端 GET 指令执行请求
func (c *Client) Get(ctx context.Context, key string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err 
	}

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})

	// 替代 io.Copy
	
	_, err = conn.Write(buf.Bytes())

	tmpBuf := make([]byte, 1024)
	n, err := conn.Read(tmpBuf)
	fmt.Println(string(tmpBuf[:n]))

	return err
}