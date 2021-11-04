package main

import (
	"VatIdValidator/internal/http/rest"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(handler *rest.Handler, cfg *MainConfig) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "not found")
	})

	//single routes
	r.GET("/health", handler.Health)

	//group routes
	company := r.Group("/validator")
	{
		company.GET("/vat/:id", handler.VatValidator)
	}

	//sample auth on one of domains
	if cfg.Server.AuthEnabled {
		company.Use(gin.BasicAuth(gin.Accounts{
			cfg.Server.User: cfg.Server.Pass,
		}))
	}

	var AllowedRoutes = make(map[string]bool, 0)

	routes := r.Routes()
	for _, i := range routes {
		AllowedRoutes[i.Path] = true
	}

	return r
}
