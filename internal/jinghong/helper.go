package jinghong

import (
	"jonghong/internal/jinghong/store"
	emailservice "jonghong/internal/pkg/emailservicee"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/log"
	"jonghong/internal/pkg/model"
	"jonghong/pkg/db"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// 各种配置的初始化
// 1. viper配置文件
// 2. 从配置文件中读取log配置
// 3. 从配置文件中读取数据库配置

const (
	defaultConfigName = "default.yaml"
)

func initDirPath() {
	home, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	known.SetHomeDir(filepath.Join(filepath.Dir(home), ".."))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configPath := filepath.Join(known.HomeDir, "/configs", defaultConfigName)
		viper.SetConfigFile(configPath)
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {

	}
}

func logOptions() *log.Options {

	lp := &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}

	return lp
}

func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}
	ins, err := db.NewMySql(dbOptions)

	if err != nil {
		return err
	}

	// 检查数据库内是否存在目标表
	ins.AutoMigrate(&model.UserM{})

	_ = store.NewStore(ins)

	return nil

}

func initService() error {
	return emailservice.InitMailService()
}
