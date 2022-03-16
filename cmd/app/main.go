package main

import (
	"context"
	"example/internal/config"
	"example/internal/controller"
	"example/internal/http"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sendya/pkg/env"
	"go.uber.org/fx"
)

func main() {
	flag.Parse()

	// setter timezone
	os.Setenv("TZ", "Asia/Shanghai")
	cst := time.FixedZone("CST", 8*3600)
	time.Local = cst

	ctx := context.Background()

	if app := setupApp(ctx); app != nil {
		fmt.Println("Serve running.")

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		_ = app.Stop(ctx)

		fmt.Println("\r\nBye.")
	}
}

func setupApp(ctx context.Context) *fx.App {
	app := fx.New(
		// if need provide log, you can remove `fx.NopLogger`.
		fx.NopLogger,
		// provide
		fx.Options(
			fx.Provide(env.CompileInfo),
			fx.Provide(config.New),

			fx.Provide(http.New),
		),

		// inject
		fx.Options(
			// handle controllers
			controller.Modules,
		),
	)

	if err := app.Start(ctx); err != nil {
		log.Fatal("app start err", err.Error())
		return nil
	}

	return app
}
