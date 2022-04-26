package rootcmd

import (
	"context"
	"example/internal/config"
	"example/internal/controller"
	"example/internal/http"
	"fmt"
	"github.com/sendya/pkg/env"
	"github.com/sendya/pkg/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
)

// project MyApp CLI
var (
	rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "MyApp is a fast website and server uptime monitoring.",
		Long: `A fast website and server uptime monitoring.
Complete documentation is available at https://yoursite.com`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of MyApp",
		Long:  "All software has versions. This is MyApp",
		Run: func(cmd *cobra.Command, args []string) {
			env.CompileInfo().Print("MyApp(rootCmd)")
			os.Exit(0)
		},
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "serve MyApp webserver",
		Long:  "run and handler http webserver in MyApp",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			if app := setupApp(ctx); app != nil {
				fmt.Println("Serve running.")

				ch := make(chan os.Signal, 1)
				signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
				<-ch
				_ = app.Stop(ctx)

				fmt.Println("\r\nBye.")
			}
		},
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize project default db",
		Long:  "initialize MyApp database schema",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("init...")
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)

	rootCmd.PersistentFlags().StringVar(&config.AppEnv, "env", "prod", "dev, test, prod")
	rootCmd.PersistentFlags().StringVar(&config.DefFileName, "config", "config", "config filename")
}

func Execute() error {
	return rootCmd.Execute()
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
		log.Fatal("app start err", log.ErrorField(err))
		return nil
	}

	return app
}
