package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
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

	// Put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	str := `[{"path":"D:/Code/logs/test.log","topic":"web_log"},{"path":"F:/Code/logs/test.log","topic":"mysql_log"}]`
	_, err = cli.Put(ctx, "collect_log_conf", str)
	if err != nil {
		fmt.Printf("Put to etcd failed, err: %v", err)
		return
	}
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	gr, err := cli.Get(ctx, "collect_log_conf")
	if err != nil {
		fmt.Printf("Get from etcd failed, err: %v", err)
		return
	}
	cancel()

	for _, ev := range gr.Kvs {
		fmt.Printf("key: %s value: %s\n", ev.Key, ev.Value)
	}
}
