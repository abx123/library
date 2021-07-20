package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"

	"library/handler"
	"library/handler/middleware"
)

type router struct {
	port int
	conn *sqlx.DB
}

func NewRouter(port int, conn *sqlx.DB) *router {
	return &router{
		port: port,
		conn: conn,
	}
}

func (router *router) InitRouter() *echo.Echo {
	handler := handler.NewHandler(router.conn)
	r := echo.New()

	// Middleware
	r.Pre(middleware.Middleware())

	p := prometheus.NewPrometheus("library", nil)
	p.Use(r)

	r.GET("/ping", handler.Ping)
	r.GET("/:userId/book/:isbn", handler.GetBook)
	r.GET("/:userId/books", handler.ListBook)
	r.POST("/:userId/book/:isbn", handler.UpsertBook)
	r.GET("/book/:isbn", handler.GetNewBook)

	r.Start(fmt.Sprintf(":%d", router.port))
	return r
}
