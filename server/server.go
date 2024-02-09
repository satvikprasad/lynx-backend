package server

import (
	"fmt"
	"lynx-backend/db"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Context struct {
	Database db.DB

	context *gin.Context
}

type Server struct {
	r  *gin.Engine
	db db.DB

	port string
}

func (c *Context) WriteJSON(code int, v interface{}) {
	c.context.Header("Content-Type", "application/json")
	c.context.JSON(code, v)
}

func (c *Context) BindJSON(v interface{}) error {
	return c.context.BindJSON(v)
}

func (c *Context) GetParam(key string) string {
	return c.context.Param(key)
}

func (c *Context) WriteString(code int, v string) {
	c.context.String(code, v)
}

func (c *Context) WriteHTML(code int, v string, obj any) {
	c.context.HTML(code, v, obj)
}

func CreateServer(db db.DB, port string) *Server {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static/assets", "./assets")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	return &Server{
		r:    r,
		db:   db,
		port: port,
	}
}

func (s *Server) Get(route string, fn ApiFunc) {
	s.r.GET(route, makeApiFunc(s.db, fn))
}

func (s *Server) Post(route string, fn ApiFunc) {
	s.r.POST(route, makeApiFunc(s.db, fn))
}

func (s *Server) Listen() {
	s.r.Run(":" + s.port)
}

func makeApiFunc(db db.DB, fn ApiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "db is nil"})
			return
		}

		ctx := &Context{
			Database: db,
			context:  c,
		}

		if err := fn(ctx); err != nil {
			fmt.Printf("Error: %v", err)
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

type ApiFunc func(c *Context) error
