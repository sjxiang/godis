package main

import (
	"io"
	"net"
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
	buf := make([]byte, 1024)

	for {
		// input
		n, err := p.Receive(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
	
		// 
		p.msgCh <- Message{
			data: msgBuf,
			peer: p,
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