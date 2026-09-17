package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/link"
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

	router.Run(app.addr)
}
