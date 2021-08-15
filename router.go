package main

import (
	"fmt"

	goisbn "github.com/abx123/go-isbn"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"

	"github.com/abx123/library/handler"
	"github.com/abx123/library/handler/middleware"
	"github.com/abx123/library/repo"
	"github.com/abx123/library/services"
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

	dbRepo := repo.NewDbRepo(router.conn)
	gi := goisbn.NewGoISBN(goisbn.DEFAULT_PROVIDERS)
	handler := handler.NewHandler(services.NewDbService(dbRepo), services.NewBookService(gi))
	r := echo.New()

	// Middleware
	r.Pre(middleware.Middleware())

	p := prometheus.NewPrometheus("library", nil)
	p.Use(r)

	r.GET("/ping", handler.Ping)
	r.GET("/:userId/book/:isbn", handler.GetBook)
	r.GET("/:userId/books", handler.ListBook)
	r.POST("/:userId/book", handler.UpsertBook)
	r.GET("/book/:isbn", handler.GetNewBook)

	r.Start(fmt.Sprintf(":%d", router.port))
	return r
}
