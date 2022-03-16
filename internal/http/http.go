package http

import (
	"example/internal/config"

	"github.com/gin-gonic/gin"
)

func New(conf *config.Config) (*gin.Engine, *gin.RouterGroup) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// add middleware
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")

	// v1.Use(authority.Def)

	return r, v1
}
