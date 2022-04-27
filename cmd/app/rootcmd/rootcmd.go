package rootcmd

import (
	"context"
	"example/internal/core/config"
	"example/internal/core/db"
	"example/internal/service"
	"fmt"
	"github.com/sendya/pkg/env"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"example/internal/controller"
	"example/internal/http"
)

type ReadyApp = struct{}

// project MyApp CLI
var (
	rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "MyApp is a fast website and server uptime monitoring.",
		Long: `A fast website and server uptime monitoring.
Complete documentation is available at https://yoursite.com`,
		Run: func(cmd *cobra.Command, args []string) {
			exec := executable()
			fmt.Printf("you can use some command to run this app. e.g. %s serve\n", exec)

			runServe()
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
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "generate project config.yml file",
		Long:  "Automatically generate project config.yml file.",
		Run: func(cmd *cobra.Command, args []string) {
			config.Genconfig = true
			config.New()
		},
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "serve MyApp webserver",
		Long:  "run and handler http webserver in MyApp",
		Run: func(cmd *cobra.Command, args []string) {
			runServe()
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
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().StringVar(&config.AppEnv, "env", "prod", "dev, test, prod")
	rootCmd.PersistentFlags().StringVar(&config.DefFileName, "config", "config", "config filename")
}

func executable() string {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	return exec
}

func Execute() error {
	return rootCmd.Execute()
}

func runServe() {
	ready := make(chan ReadyApp, 1)
	ctx := context.WithValue(context.Background(), "-15", ready)

	if app := setupApp(ctx); app != nil {
		<-ready

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
			fx.Provide(db.New),

			fx.Provide(http.New),
		),

		// services
		fx.Options(
			fx.Provide(service.NewUserService),
		),

		// inject
		fx.Options(
			// handle controllers
			controller.Modules,
		),
	)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	return app
}
