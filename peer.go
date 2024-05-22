package main

import (
	"fmt"
	"io"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/resp"
)


type Peer struct {
	conn  net.Conn
	msgCh chan Message
}

func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
	
		v, _, err := rd.ReadValue()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("resp read value error")  // 识别出错
		}

		var cmd Command
		if v.Type() == resp.Array {
			rawCmd := v.Array()[0]

			switch rawCmd.String() {
			case "SET":
				if len(v.Array()) != 3 {
					return fmt.Errorf("invalid number of variables for SET command")
				}
				cmd = &SetCommand{
					key:   v.Array()[1].Bytes(),
					value: v.Array()[2].Bytes(),
				}
			case "GET":
				if len(v.Array()) != 2 {
					return fmt.Errorf("invalid number of variables for GET command")
				}
				cmd = &GetCommand{
					key: v.Array()[1].Bytes(),
				}
			case "HELLO":
				cmd = &HelloCommand{
					value: v.Array()[1].String(),
				}
			default:
				fmt.Println("got this unhandled command", rawCmd)
			}

			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,	
			}
		}
	}

	return nil
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) Receive(msg []byte) (int, error) {
	return p.conn.Read(msg)
}

