package tailf

import (
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/etcd"
	"liujun/Time_ELK/kafka"
)

var (
	err error
)

type TailTask struct {
	Path  string     `json:"path"`
	Topic string     `json:"topic"`
	Tail  *tail.Tail `json:"tail"`
}

func NewTailf(msg etcd.EtcdMsg) *TailTask {
	tail_task := &TailTask{
		Path:  msg.Path,
		Topic: msg.Topic,
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
		line, ok := <-tt.Tail.Lines
		if ok {
			logrus.Println(line.Text)
			msg := &sarama.ProducerMessage{}
			msg.Topic = tt.Topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.Msg <- msg
		}
	}
}

func Init(msgs []etcd.EtcdMsg) error {
	for _, msg := range msgs {
		tail_task := NewTailf(msg)
		tail_task, err = tail_task.Init() //根据etcd创建任务
		if err != nil {
			return err
		}
		//启动任务
		go tail_task.Run()
	}
	return nil
}
