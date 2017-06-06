package server

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/inl2_das/conf"
	"github.com/giskook/inl2_das/conn"
	"log"
	"net"
	"time"
)

type ServerConfig struct {
	Listener      *net.TCPListener
	AcceptTimeout time.Duration
}

type Server struct {
	config           *ServerConfig
	srv              *gotcp.Server
	checkconnsticker *time.Ticker
}

var Gserver *Server

func SetServer(server *Server) {
	Gserver = server
}

func GetServer() *Server {
	return Gserver
}

func NewServer(srv *gotcp.Server, config *ServerConfig) *Server {
	serverstatistics := conf.GetConf().Server.ServerStatistics
	return &Server{
		config:           config,
		srv:              srv,
		checkconnsticker: time.NewTicker(time.Duration(serverstatistics) * time.Second),
	}
}

func (s *Server) Start() {
	go s.checkStatistics()

	s.srv.Start(s.config.Listener, s.config.AcceptTimeout)
}

func (s *Server) Stop() {
	s.srv.Stop()
	s.checkconnsticker.Stop()

}

func (s *Server) checkStatistics() {
	for {
		<-s.checkconnsticker.C
		log.Printf("---------------------Total Connections : %d---------------------\n", conn.NewConns().GetCount())
	}
}
