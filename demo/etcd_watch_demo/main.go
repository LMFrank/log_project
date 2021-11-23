package main

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

import (
	"fmt"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v", err)
		return
	}

	defer cli.Close()

	// Watch
	watchCh := cli.Watch(context.Background(), "name")
	for wresp := range watchCh {
		for _, evt := range wresp.Events {
			fmt.Printf("type: %s, key: %s, value: %s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
		}
	}

}
