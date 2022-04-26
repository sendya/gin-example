package main

import (
	"example/cmd/app/rootcmd"
	"fmt"
	"os"
	"time"
)

// @title                       Uncrash Core
// @version                     1.0.0
// @description                 Fast website and server uptime monitoring.
// @BasePath                    /api
// @securityDefinitions.apikey  Authorization
// @in                          header
// @name                        Authorization
func main() {
	// setter timezone
	os.Setenv("TZ", "Asia/Shanghai")
	cst := time.FixedZone("CST", 8*3600)
	time.Local = cst

	if err := rootcmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Fprintln: %v\n", err)
		os.Exit(1)
	}
}
