package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"log_project/logagent/kafka"
	"log_project/logagent/tailfile"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

func main() {
	var configObj = new(Config)

	// 0. 读取配置文件 go-ini
	/*
		cfg, err := ini.Load("./conf/config.ini")
		if err != nil {
			logrus.Error("load config failed, err: %v", err)
			return
		}
		kafkaAddr := cfg.Section("kafka").Key("address").String()
		fmt.Println(kafkaAddr)
	*/

	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err: %v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)

	// 1. 初始化连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address})
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	// 2. 根据配置中的日志路径初始化tail
	err = tailfile.Init(configObj.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Errorf("init tailfile failed, err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")

	// 3. 把日志通过sarama发往kafka
}
