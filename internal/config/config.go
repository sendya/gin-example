package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sendya/pkg/log"
	"github.com/spf13/viper"
)

var (
	AppEnv       string = "prod"
	defFileName  string = "config"
	defFileExt   string = "yml"
	GlobalCofnig *Config
)

type Config struct {
	APPName string `mapstructure:"app-name" json:"app-name"`
	Host    string `mapstructure:"host" json:"host"`
	Port    int    `mapstructure:"port" json:"port"`

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

func init() {
	flag.StringVar(&AppEnv, "env", "prod", "dev, test, prod")
}

func New() *Config {
	var (
		err  error
		conf *Config
	)

	viper.SetConfigName(fmt.Sprintf("%s.%s", defFileName, AppEnv))
	viper.SetConfigType(defFileExt)
	// viper.AddConfigPath("/etc/example/")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		log.Warn("load config err", log.String("msg", err.Error()))
		viper.SetConfigName(defFileName)
		log.Info("try load default config", log.String("filename", defFileName))
		if err = viper.ReadInConfig(); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	}

	if err = viper.Unmarshal(&conf); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}

	log.Info("conf loaded", log.String("level", conf.Logger.Level), log.String("env", AppEnv))
	log.SetLevel(conf.Logger.Level)

	if AppEnv == "prod" {
		fileLogger(conf)
	}

	GlobalCofnig = conf

	return conf
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
