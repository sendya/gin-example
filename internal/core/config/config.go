package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sendya/pkg/log"
	"github.com/spf13/viper"
)

var (
	AppEnv       string = "prod"
	DefFileName  string = "config"
	DefFileExt   string = "yml"
	Genconfig    bool
	GlobalConfig *Config
)

type Config struct {
	App struct {
		Name string `mapstructure:"name" json:"name"`
		Host string `mapstructure:"host" json:"host"`
		Port int    `mapstructure:"port" json:"port"`
	} `mapstructure:"app" json:"app"`

	Database struct {
		Type     string   `mapstructure:"type" json:"type"`
		DSN      []string `mapstructure:"dsn" json:"dsn"`
		Migrator bool     `mapstructure:"migrator" json:"migrator"`
		Debug    bool     `mapstructure:"debug" json:"debug"`
	} `mapstructure:"database" json:"database"`

	Redis struct {
		Addr     string `mapstructure:"addr" json:"addr"`
		Password string `mapstructure:"password" json:"password"`
		DB       uint   `mapstructure:"db" json:"db"`
	} `mapstructure:"redis" json:"redis"`

	Logger struct {
		Level  string `mapstructure:"level" json:"level"`
		Path   string `mapstructure:"path" json:"path"`
		Caller bool   `mapstructure:"caller" json:"caller"`
	} `mapstructure:"logger" json:"logger"`
}

func New() (*Config, *viper.Viper) {
	var (
		err  error
		conf *Config
		v    *viper.Viper
	)
	cf := envFile(DefFileName, AppEnv)
	v = viper.New()

	v.SetConfigName(cf)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.AddConfigPath(".")

	// set default config
	// base app
	v.SetDefault("app.name", "myapp")
	v.SetDefault("app.host", "localhost")
	v.SetDefault("app.port", 9000)
	// logger
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.path", "./logs")
	v.SetDefault("logger.caller", true)
	// database
	v.SetDefault("database.type", "mysql")
	v.SetDefault("database.dsn", []string{"root:root@tcp(127.0.0.1:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local"})
	v.SetDefault("database.debug", true)
	v.SetDefault("database.migrator", true)
	// redis
	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 1)

	if Genconfig {
		filename := filepath.Join(filepath.Dir("config/"), cf)
		v.WriteConfigAs(filename)
		fmt.Println("writter config " + filename)
		os.Exit(0)
	}

	log.Info("load env config", log.String("filename", cf))
	if err = v.ReadInConfig(); err != nil {
		log.Warn("load env config not found")

		cf = envFile(DefFileName, "")
		v.SetConfigName(cf)
		log.Info("load def config", log.String("filename", cf))
		if err = v.ReadInConfig(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	}

	if err = v.Unmarshal(&conf); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}

	log.Info("conf load", log.String("level", conf.Logger.Level), log.String("env", AppEnv))
	log.SetLevel(conf.Logger.Level)

	if AppEnv == "prod" {
		fileLogger(conf)
	}

	GlobalConfig = conf
	return conf, v
}

func fileLogger(conf *Config) {
	var tops = []log.TeeOption{
		{
			Filename: filepath.Join(conf.Logger.Path, "access.log"),
			Ropt: log.RotateOptions{
				MaxSize:    10,
				MaxAge:     1,
				MaxBackups: 6,
				Compress:   true,
			},
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
		{
			Filename: filepath.Join(conf.Logger.Path, "error.log"),
			Ropt: log.RotateOptions{
				MaxSize:    10,
				MaxAge:     1,
				MaxBackups: 6,
				Compress:   true,
			},
			Lef: func(lvl log.Level) bool {
				return lvl > log.InfoLevel
			},
		},
	}

	opts := make([]log.Option, 0, 2)
	if conf.Logger.Caller {
		opts = append(opts, log.WithCaller(true), log.AddCallerSkip(0))
	}

	rotate := log.NewTeeWithRotate(tops, opts...)
	log.ResetDefault(rotate)
}

func envFile(filename string, env string) string {
	if env == "" {
		return fmt.Sprintf("%s.%s", filename, DefFileExt)
	}
	return fmt.Sprintf("%s.%s.%s", filename, env, DefFileExt)
}
