package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	Client *clientv3.Client
	err    error
)

func Init(addr []string) error {
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Second * 5,
	})
	return err
}

type EtcdMsg struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

//topic:path
func GetEtcdValue(etcd_key string) (data []EtcdMsg, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := Client.Get(ctx, etcd_key)
	defer cancel()
	if err != nil {
		return nil, nil
	}
	if len(resp.Kvs) <= 0 {
		return nil, nil
	}
	res := resp.Kvs[0]
	err = json.Unmarshal(res.Value, &data)
	if err != nil {
		return
	}
	return
}
