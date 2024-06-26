package main

import (
	"net"
	"os"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/resp"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})  // 控制台输出、文件输出	
}

const (
	defaultListenAddr = ":5001"
)


type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	msgCh     chan Message
	quitCh    chan struct{}
	kv        *KV
}

type Message struct {
	cmd  Command
	peer *Peer
}



func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	log.Info().Str("监听地址", s.ListenAddr).Msg("服务器正在运行...")
	
	return s.acceptLoop()
}

// 主进程
func (s *Server) loop() {
	for {
		select {
		case rawMsg := <- s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				log.Error().Err(err).Str("监听地址", s.ListenAddr).Msg("raw message error")
			}
		case <- s.quitCh:
			return
		case peer := <- s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

// 处理请求
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error().Err(err).Msg("accept error")
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	
	log.Info().Str("远程地址", conn.RemoteAddr().String()).Msg("new peer connected")

	if err := peer.readLoop(); err != nil {
		log.Error().Err(err).Str("远程地址", conn.RemoteAddr().String()).Msg("peer receive err")
	}
}


func (s *Server) handleRawMessage(msg Message) error {
	
	switch v := msg.cmd.(type) {
	case *SetCommand:
		if err := s.kv.Set(v.key, v.value); err != nil {
			return err
		}
		if err := resp.NewWriter(msg.peer.conn).WriteString("OK"); err != nil {
			return err
		}
	case *GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		if err := resp.NewWriter(msg.peer.conn).WriteString(string(val)); err != nil {
			return err
		}
	case *HelloCommand:
		spec := map[string]string{
			"server": "reids",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			return fmt.Errorf("peer send error: %s", err)
		}
	}

	return nil 
}




func main() {
	cfg := Config{}
	srv := NewServer(cfg)
	
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("server start error")
	}
}