package http

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sendya/pkg/log"
	"go.uber.org/fx"

	"example/internal/config"
)

func New(
	lc fx.Lifecycle,
	conf *config.Config,
) (*gin.Engine, *gin.RouterGroup) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// add middleware
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")

	// v1.Use(authority.Def)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Infof("listening web serve http://localhost:%d", conf.Port)
				log.Infof("listening swagserve http://localhost:%d/swagger/index.html", conf.Port)

				if err := r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
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
