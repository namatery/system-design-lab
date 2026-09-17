package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/link"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/redirect"
)

type Application struct {
	db   *gocql.Session
	addr string
}

func (app *Application) mount() {
	router := gin.Default()

	// Links
	{
		service := link.NewLinkService(app.db)
		handler := link.NewLinkController(service)

		v1 := router.Group("/v1/link")
		v1.POST("", handler.CreateShortenLink)
	}

	// Redirect
	{
		service := redirect.NewRedirectService(app.db)
		handler := redirect.NewRedirectController(service)

		v1 := router.Group("/")
		v1.GET(":alias", handler.Redirect)
	}

	router.Run(app.addr)
}
