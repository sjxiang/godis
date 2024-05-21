package main

import (
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	msgCh     chan []byte
	quitCh    chan struct{}
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
		msgCh:     make(chan []byte),
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
		log.Error().Err(err).Str("远程地址", conn.RemoteAddr().String()).Msg("peer read err")
	}
}


func (s *Server) handleRawMessage(rawMsg []byte) error {
	LogMessage(Blue, string(rawMsg))
	return nil 
}

func main() {
	srv := NewServer(Config{})
	
	if err := srv.Start(); err != nil {
		panic(err)
	}
}