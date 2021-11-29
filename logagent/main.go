package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"log_project/logagent/common"
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
	ip, err := common.GetOutboundIP()
	if err != nil {
		logrus.Errorf("get ip failed, err: %v", err)
		return
	}

	var configObj = new(Config)

	// 读取配置文件 go-ini
	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err: %v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)

	// 初始化kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err: %v", err)
		return
	}
	logrus.Info("init kafka success!")

	// 初始化etcd
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed, err: %v", err)
		return
	}

	// 从etcd中拉取要收集日志的配置项
	collectKey := fmt.Sprintf(configObj.EtcdConfig.CollectKey, ip)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err: %v", err)
		return
	}

	fmt.Printf("allConf: %s\n", allConf)

	// 监控etcd中 configObj.EtcdConfig.CollectKey 对应值的变化
	go etcd.WatchConf(collectKey)

	// 根据配置中的日志路径初始化tail，把日志通过sarama发往kafka
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tailfile failed, err: %v", err)
		return
	}
	logrus.Info("init tailfile success!")

	select {}
}
