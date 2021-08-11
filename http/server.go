package http

import (
	"context"
	"encoding/base64"
	"errors"
	"expinc/sunagent/common"
	"expinc/sunagent/http/handlers"
	"expinc/sunagent/log"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const urlPrefix = "/api/v1"

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type stopInfo struct {
	waitSec int64
	err     error
}

type Server struct {
	ip   string
	port uint

	engine *gin.Engine

	// channel element values as gracefully stopping seconds
	// after the gracefully stopping seconds, server will be stopped gracefully
	quit chan int64
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

	// register handlers
	server.registerHandlers()

	// starting HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", server.ip, server.port),
		Handler: server.engine,
	}
	server.quit = make(chan int64, 1)
	srvErr := make(chan error, 1)
	go func() {
		var err error
		if err = srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Info("HTTP server stopped")
		}

		if err == http.ErrServerClosed {
			err = nil
		}
		srvErr <- err
	}()

	// stopping HTTP server
	waitSec := <-server.quit
	log.Info("Stopping HTTP server...")
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Duration(waitSec)*time.Second)
	defer cancel()
	if err := srv.Shutdown(stopCtx); err != nil {
		log.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	return <-srvErr
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
	// management
	server.engine.GET(urlPrefix+"/info", handlerProxy(handlers.GetInfo))
	server.engine.POST(urlPrefix+"/terminate", handlerProxy(server.terminate))

	// file
	server.engine.GET(urlPrefix+"/fileMeta", handlerProxy(handlers.GetFileMeta))
	server.engine.GET(urlPrefix+"/file", handlerProxy(handlers.GetFileContent))
	server.engine.POST(urlPrefix+"/file", handlerProxy(handlers.CreateFile))
	server.engine.PUT(urlPrefix+"/file", handlerProxy(handlers.OverwriteFile))
	server.engine.DELETE(urlPrefix+"/file", handlerProxy(handlers.DeleteFile))

	// processes
	server.engine.GET(urlPrefix+"/processes/:pidOrName", handlerProxy(handlers.GetProcInfo))
	server.engine.POST(urlPrefix+"/processes/:pidOrName/kill", handlerProxy(handlers.KillProc))
	server.engine.POST(urlPrefix+"/processes/:pidOrName/terminate", handlerProxy(handlers.TermProc))
}

func (server *Server) terminate(ctx *gin.Context) {
	handlers.RespondSuccessfulJson(ctx, http.StatusNoContent, nil)
	waitSecVals, ok := ctx.Request.URL.Query()["waitSec"]

	// wait for 3 seconds for gracefully stop by default
	waitSec := int64(3)
	if ok {
		var err error
		waitSec, err = strconv.ParseInt(waitSecVals[0], 10, 64)
		if nil != err {
			waitSec = 0
		}
	}
	server.quit <- waitSec
}
