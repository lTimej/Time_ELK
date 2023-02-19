package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"liujun/Time_ELK/common"
	"liujun/Time_ELK/tail_task_mgr"
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

//topic:path
func GetEtcdValue(etcd_key string) (data []common.EtcdMsg, err error) {
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

func Watch(etcd_key string) {
	ew := Client.Watch(context.Background(), etcd_key)
	for wresp := range ew {
		for _, ev := range wresp.Events {
			data := []common.EtcdMsg{}
			err := json.Unmarshal(ev.Kv.Value, &data)
			if err != nil {
				return
			}
			tail_task_mgr.PutNewChan(data)
		}
	}
}
