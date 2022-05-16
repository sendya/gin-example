package main

import (
	"os"
	"time"

	"example/cmd/app/rootcmd"
)

// @title                       Gin-Example
// @version                     1.0.0
// @description                 Fast website and server monitoring.
// @BasePath                    /api
// @securityDefinitions.apikey  Authorization
// @in                          header
// @name                        Authorization
func main() {
	// setter timezone
	_ = os.Setenv("TZ", "Asia/Shanghai")
	cst := time.FixedZone("CST", 8*3600)
	time.Local = cst

	if err := rootcmd.Execute(); err != nil {
		panic(err)
	}
}
