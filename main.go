package main

import (
	"crud-api/handler"
	"crud-api/middleware"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
}

var (
	ServiceName = "crud-api"
	port        = "10000"
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
}

func NewServer() *Server {
	r := gin.New()
	return &Server{
		router: r,
	}
}

func main() {

	server := NewServer()
	server.Router()

	zapLin, _ = zap.NewProduction()

	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: server.router,
	}

	err := server.httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (s *Server) Router() {

	grouter := s.router.Group("v1")
	grouter.Use(GinRecovery(true))
	grouter.Use(middleware.Cors())
	// grouter.Use(middleware.CorsByRules())
	// demo
	grouter.GET("/set", handler.Set())
	grouter.GET("/get-a", handler.Get())
	grouter.GET("/get-b", handler.GetB())
	grouter.GET("/health", handler.Health())

}

var zapLin *zap.Logger

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					fmt.Println("3333")

					zapLin.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zapLin.Info("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zapLin.Info("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
