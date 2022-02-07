package srv

import (
	"github.com/bugfan/de"
	"github.com/bugfan/srv"
	"github.com/bugfan/trojan-auth/env"
	"github.com/sirupsen/logrus"
)

func init() {
	de.SetKey(env.Get("des_key"))
}

func NewServer(addr string) *Server {

	s := srv.New(addr)
	s.Handle("/wireguard", &Wireguard{}) // set wg handler
	s.Handle("/config", &Config{})

	return &Server{
		addr,
		s,
	}
}

type Server struct {
	addr string
	*srv.Server
}

func (s *Server) Run() {
	logrus.Fatal(s.Server.Run())
}
