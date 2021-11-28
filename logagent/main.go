package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"log_project/logagent/etcd"
	"log_project/logagent/kafka"
	"log_project/logagent/tailfile"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func main() {
	var configObj = new(Config)

	// 0. 读取配置文件 go-ini
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err: %v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)

	// 1. 初始化连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err: %v", err)
		return
	}
	logrus.Info("init kafka success!")

	// 2. 初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed, err: %v", err)
		return
	}

	// 3. 从etcd中拉取要收集日志的配置项
	allConf, err := etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err: %v", err)
		return
	}

	fmt.Printf("allConf: %s\n", allConf)

	// 4. 派一个小弟去监控etcd中 configObj.EtcdConfig.CollectKey 对应值的变化
	go etcd.WatchConf(configObj.EtcdConfig.CollectKey)

	// 5`. 根据配置中的日志路径初始化tail，把日志通过sarama发往kafka
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tailfile failed, err: %v", err)
		return
	}
	logrus.Info("init tailfile success!")

	select {}
}
