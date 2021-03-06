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
	"strings"
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
	ip          string
	port        uint
	authMethod  string
	authCred    interface{}
	tlsCertFile string
	tlsKeyFile  string

	engine         *gin.Engine
	authMiddleware gin.HandlerFunc

	// channel element values as gracefully stopping seconds
	// after the gracefully stopping seconds, server will be stopped gracefully
	quit chan int64
}

// FIXME: move to a package of security
type BasicAuthCred struct {
	User     string
	Password string
}

func New(ip string, port uint, authMethod string, authCred interface{}, certFile, keyFile string) *Server {
	return &Server{
		ip:          ip,
		port:        port,
		authMethod:  authMethod,
		authCred:    authCred,
		tlsCertFile: certFile,
		tlsKeyFile:  keyFile,
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
		handlers.RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	}
	server.engine.Use(gin.CustomRecovery(recoverFunc))

	// register handlers
	err := server.registerHandlers()
	if nil != err {
		return err
	}

	// starting HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", server.ip, server.port),
		Handler: server.engine,
	}
	server.quit = make(chan int64, 1)
	srvErr := make(chan error, 1)
	go func() {
		var err error
		if "" != strings.TrimSpace(server.tlsCertFile) && "" != strings.TrimSpace(server.tlsKeyFile) {
			err = srv.ListenAndServeTLS(server.tlsCertFile, server.tlsKeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && errors.Is(err, http.ErrServerClosed) {
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

func (server *Server) initAuth() error {
	methodStr := strings.TrimSpace(server.authMethod)
	if "" == methodStr || "none" == methodStr {
		return nil
	}

	if "basic" == methodStr {
		cred, ok := server.authCred.(BasicAuthCred)
		invalidAuthCredMsg := "Invalid authentication credential type"
		if !ok {
			return common.NewError(common.ErrorInvalidParameter, invalidAuthCredMsg)
		}
		if "" == strings.TrimSpace(cred.User) || "" == strings.TrimSpace(cred.Password) {
			return common.NewError(common.ErrorInvalidParameter, invalidAuthCredMsg)
		}

		server.authMiddleware = gin.BasicAuth(gin.Accounts{
			cred.User: cred.Password,
		})
	} else {
		msg := fmt.Sprintf("Invalid authentication method \"%s\"", server.authMethod)
		return common.NewError(common.ErrorInvalidParameter, msg)
	}

	return nil
}

func (server *Server) registerHandlers() error {
	if err := server.initAuth(); nil != err {
		return err
	}
	var routerGroup *gin.RouterGroup
	routerGroup = &server.engine.RouterGroup
	middlewares := make([]gin.HandlerFunc, 0)
	if nil != server.authMiddleware {
		middlewares = append(middlewares, server.authMiddleware)
	}

	// management
	routerGroup.GET(urlPrefix+"/info", handlerProxy(handlers.GetInfo))
	routerGroup = server.engine.Group(urlPrefix+"/terminate", middlewares...)
	routerGroup.POST("", handlerProxy(server.terminate))

	// file
	routerGroup = server.engine.Group(urlPrefix+"/file", middlewares...)
	routerGroup.GET("/meta", handlerProxy(handlers.GetFileMeta))
	routerGroup.GET("", handlerProxy(handlers.GetFileContent))
	routerGroup.POST("", handlerProxy(handlers.CreateFile))
	routerGroup.PUT("", handlerProxy(handlers.OverwriteFile))
	routerGroup.POST("/append", handlerProxy(handlers.AppendFile))
	routerGroup.DELETE("", handlerProxy(handlers.DeleteFile))

	// processes
	routerGroup = server.engine.Group(urlPrefix+"/processes", middlewares...)
	routerGroup.GET("/:pidOrName", handlerProxy(handlers.GetProcInfo))
	routerGroup.POST("/:pidOrName/kill", handlerProxy(handlers.KillProc))
	routerGroup.POST("/:pidOrName/terminate", handlerProxy(handlers.TermProc))

	// system
	routerGroup = server.engine.Group(urlPrefix+"/sys", middlewares...)
	routerGroup.GET("/info", handlerProxy(handlers.GetNodeInfo))
	routerGroup.GET("/cpus/info", handlerProxy(handlers.GetCpuInfo))
	routerGroup.GET("/cpus/stats", handlerProxy(handlers.GetCpuStat))
	routerGroup.GET("/mem/stats", handlerProxy(handlers.GetMemStat))
	routerGroup.GET("/disks/stats", handlerProxy(handlers.GetDiskInfo))
	routerGroup.GET("/net/info", handlerProxy(handlers.GetNetInfo))

	// script
	routerGroup = server.engine.Group(urlPrefix+"/script", middlewares...)
	routerGroup.POST("/execute", handlerProxy(handlers.ExecScript))

	// package
	routerGroup = server.engine.Group(urlPrefix+"/package", middlewares...)
	routerGroup.GET("/:name", handlerProxy(handlers.GetPackageInfo))
	routerGroup.POST("/:name", handlerProxy(handlers.InstallPackage)) // install by name
	routerGroup.POST("", handlerProxy(handlers.InstallPackage))       // install by file
	routerGroup.PUT("/:name", handlerProxy(handlers.UpgradePackage))  // upgrade by name
	routerGroup.PUT("", handlerProxy(handlers.UpgradePackage))        // upgrade by file
	routerGroup.DELETE("/:name", handlerProxy(handlers.UninstallPackage))

	// job
	routerGroup = server.engine.Group(urlPrefix+"/jobs", middlewares...)
	routerGroup.GET("/:id", handlerProxy(handlers.GetJobInfo))
	routerGroup.GET("", handlerProxy(handlers.GetAllJobInfo))
	routerGroup.POST("/:id/cancel", handlerProxy(handlers.CancelJob))

	// grimoire
	routerGroup = server.engine.Group(urlPrefix+"/grimoires", middlewares...)
	routerGroup.GET("/:osType", handlerProxy(handlers.GetGrimoire))
	routerGroup.POST("/:osType/arcanes/:arcaneName/cast", handlerProxy(handlers.CastArcane))
	routerGroup.PUT("/:osType/arcanes/:arcaneName", handlerProxy(handlers.SetArcane))
	routerGroup.DELETE("/:osType/arcanes/:arcaneName", handlerProxy(handlers.RemoveArcane))

	return nil
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
