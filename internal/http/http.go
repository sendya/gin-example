package http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sendya/pkg/log"
	"go.uber.org/fx"

	"example/internal/config"
)

func New(
	lc fx.Lifecycle,
	conf *config.Config,
) (*gin.Engine, *gin.RouterGroup) {
	confapp := conf.App

	gin.SetMode(gin.DebugMode)
	if config.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = io.MultiWriter(ioutil.Discard)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%v %v (%v handlers)", httpMethod, absolutePath, nuHandlers)
	}
	r := gin.New()

	// add middleware
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")

	// v1.Use(authority.Def)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Infof("listening web serve http://localhost:%d", confapp.Port)

				if err := r.Run(fmt.Sprintf("%s:%d", confapp.Host, confapp.Port)); err != nil {
					log.Fatal("start web serve ", log.String("err", err.Error()))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return r, v1
}
