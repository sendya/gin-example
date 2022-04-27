package db

import (
	"example/internal/core/config"
	"github.com/sendya/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func New(c *config.Config) *gorm.DB {
	var (
		driver gorm.Dialector
		conf   = &gorm.Config{
			Logger: gormlogger.Discard,
		}
	)

	// only support MySQL.
	driver = mysql.New(mysql.Config{
		DSN:                      c.Database.DSN[0],
		DisableDatetimePrecision: true,
		DontSupportRenameIndex:   true,
	})

	if c.Database.Debug {
		conf.Logger = NewLog(log.Named("gorm"))
	}

	conn, err := gorm.Open(driver, conf)
	if err != nil {
		log.Fatalf("Cloud not connect to database %s", c.Database.DSN[0])
	}

	// multi database conn..
	l := len(c.Database.DSN) - 1
	if l > 0 {
		sources := make([]gorm.Dialector, l)
		for i := range sources {
			sources[i] = mysql.Open(c.Database.DSN[i+1])
		}

		if err := conn.Use(dbresolver.Register(dbresolver.Config{
			Sources:  sources,
			Replicas: []gorm.Dialector{},
			Policy:   dbresolver.RandomPolicy{},
		})); err != nil {
			log.Fatal("Cloud not connect to Policy", log.ErrorField(err))
		}
	}

	log.Info("Connected mysql srv.")

	return conn
}
