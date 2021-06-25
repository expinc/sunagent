package http

import (
	"encoding/base64"
	"errors"
	"expinc/sunagent/common"
	"expinc/sunagent/http/handlers"
	"expinc/sunagent/log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const urlPrefix = "/api/v1"

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Server struct {
	ip   string
	port uint

	engine *gin.Engine
}

func New(ip string, port uint) *Server {
	return &Server{
		ip:   ip,
		port: port,
	}
}

func (server *Server) Run() error {
	log.Info(fmt.Sprintf("Starting HTTP server: IP=%s, port=%d", server.ip, server.port))
	server.engine = gin.New()

	// register panic recover function
	recoverFunc := func(ctx *gin.Context, errObj interface{}) {
		var err error
		var ok bool
		if err, ok = errObj.(common.Error); !ok {
			err = errors.New(fmt.Sprintf("%v", errObj))
		}

		log.Error(err)
		handlers.RespondFailedJson(ctx, http.StatusInternalServerError, err)
	}
	server.engine.Use(gin.CustomRecovery(recoverFunc))

	server.registerHandlers()

	return server.engine.Run(fmt.Sprintf("%s:%d", server.ip, server.port))
}

func generateTraceId() string {
	uuidValue := uuid.New()
	return base64.URLEncoding.EncodeToString(uuidValue[:])
}

func handlerProxy(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get(common.TraceIdContextKey)
		if "" == traceId {
			traceId = generateTraceId()
		}
		context.Set(common.TraceIdContextKey, traceId)

		log.Info(fmt.Sprintf(
			"Start handling request: method=%s, URL=%s, traceId=%s",
			context.Request.Method,
			context.Request.URL,
			traceId))
		handler(context)
		log.Info(fmt.Sprintf("Finish handling request: method=%s, URL=%s, traceId=%s, status=%d",
			context.Request.Method,
			context.Request.URL,
			traceId,
			context.Value("status")))
	}
}

func (server *Server) registerHandlers() {
	server.engine.GET(urlPrefix+"/info", handlerProxy(handlers.GetInfo))
}
