package main

import (
	"fmt"
	"github.com/sllt/tao/core/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"os/signal"
	"reportData/device"
	"reportData/global"
	"syscall"
)

type Config struct {
	LevelUrl   string
	DataSource string
}

func InitGorm(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Error),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	var c Config
	conf.MustLoad("config.yaml", &c)

	db, err := InitGorm(c.DataSource)
	if err != nil {
		fmt.Println("初始化数据库失败：", err)
		os.Exit(1)
		return
	}
	global.DB = db

	//液位 ipPort
	if err := device.LoadDevice(c.LevelUrl); err != nil {
		fmt.Println("加载设备失败：", err)
		os.Exit(1)
		return
	}

	// 运行所有连接
	for _, device := range device.Devices {
		go device.ConnectAndRead()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	fmt.Printf("退出程序:%v\n", sig)
}
