package main

import (
	"crud-api/handler"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	// demo
	grouter.GET("/set", handler.Set())
	grouter.GET("/get-a", handler.Get())
	grouter.GET("/get-b", handler.GetB())
	grouter.GET("/health", handler.Health())

}
