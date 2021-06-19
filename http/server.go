package http

import (
	"expinc/sunagent/log"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Server struct {
	ip   string
	port uint16

	engine *gin.Engine
}

func New(ip string, port uint16) *Server {
	return &Server{
		ip:   ip,
		port: port,
	}
}

func (server *Server) Serve() error {
	log.Info(fmt.Sprintf("Starting HTTP server: IP=%s, port=%d", server.ip, server.port))
	server.engine = gin.New()
	return server.engine.Run(fmt.Sprintf("%s:%d", server.ip, server.port))
}
