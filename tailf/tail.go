package tailf

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/common"
	"liujun/Time_ELK/kafka"
	"time"
)

var (
	err     error
	NewChan chan []common.EtcdMsg
)

type TailTask struct {
	Path   string     `json:"path"`
	Topic  string     `json:"topic"`
	Tail   *tail.Tail `json:"tail"`
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewTailf(msg common.EtcdMsg) *TailTask {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	tail_task := &TailTask{
		Path:   msg.Path,
		Topic:  msg.Topic,
		Ctx:    ctx,
		Cancel: cancel,
	}
	return tail_task
}

func (tt *TailTask) Init() (*TailTask, error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tt.Tail, err = tail.TailFile(tt.Path, config)
	return tt, err
}

func (tt *TailTask) Run() {
	for {
		select {
		case <-tt.Ctx.Done():
			logrus.Println("任务关闭")
			return
		case line, ok := <-tt.Tail.Lines:
			if ok {
				logrus.Println(line.Text)
				msg := &sarama.ProducerMessage{}
				msg.Topic = tt.Topic
				msg.Value = sarama.StringEncoder(line.Text)
				kafka.Msg <- msg
			}
		}
	}
}
