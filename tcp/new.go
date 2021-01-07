package tcp

import (
	"github.com/upeven/SCache/cache"
	"net"

	"SCache/cluster"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	//监听tcp端口
	l, e := net.Listen("tcp", ":12346")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func New(c cache.Cache,n cluster.Node) *Server {
	return &Server{c,n}
}
