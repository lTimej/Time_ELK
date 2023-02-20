package es

import (
	"context"
	"fmt"
	"liujun/Time_ELK/EKB/kafka"

	elastic "github.com/olivere/elastic/v7"
)

var (
	EsClient *elastic.Client
	err      error
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}
type Msg struct {
	Info string `json:"info"`
}

func Init(addr string) error {
	EsClient, err = elastic.NewClient(elastic.SetURL(addr))
	if err != nil {
		return err
	}
	return nil
}

func Put(index string, topic string) {
	for {
		select {
		case msg := <-kafka.KfKMsg:
			m := Msg{Info: msg}
			put1, _ := EsClient.Index().
				Index(topic).
				BodyJson(m).
				Do(context.Background())
			fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		}
	}
}
