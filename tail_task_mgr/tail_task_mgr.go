package tail_task_mgr

import (
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/common"
	"liujun/Time_ELK/tailf"
)

var (
	TTM TailTaskMgr
)

type TailTaskMgr struct {
	TailTask map[string]*tailf.TailTask
	EtcdMsg  []common.EtcdMsg
	NewChan  chan []common.EtcdMsg
}

func Init(msgs []common.EtcdMsg) error {
	TTM = TailTaskMgr{
		TailTask: make(map[string]*tailf.TailTask),
		EtcdMsg:  msgs,
		NewChan:  make(chan []common.EtcdMsg),
	}
	for _, msg := range msgs {
		tail_task := tailf.NewTailf(msg)
		tail_task, err := tail_task.Init() //根据etcd创建任务
		if err != nil {
			return err
		}
		TTM.TailTask[tail_task.Path] = tail_task
		//启动任务
		go tail_task.Run()
	}
	go TTM.Watch()

	return nil
}

func (ttm *TailTaskMgr) Watch() {
	for {
		new_chans := <-ttm.NewChan
		for _, msg := range new_chans {
			if ttm.TaskExist(msg) {
				logrus.Println("任务存在")
				continue
			}
			tail_task := tailf.NewTailf(msg)
			tail_task, err := tail_task.Init() //根据etcd创建任务
			if err != nil {
				return
			}
			TTM.TailTask[tail_task.Path] = tail_task
			//启动任务
			go tail_task.Run()
		}
		for key, val := range ttm.TailTask {
			var flag bool
			for _, msg := range new_chans {
				if key == msg.Path {
					flag = true
					break
				}
			}
			if !flag {
				logrus.Println("该配置文件w不存在")
				delete(ttm.TailTask, key)
				val.Cancel()
			}
		}
	}
}

func (ttm TailTaskMgr) TaskExist(msg common.EtcdMsg) bool {
	_, ok := ttm.TailTask[msg.Path]
	return ok
}
func PutNewChan(new_chan []common.EtcdMsg) {
	TTM.NewChan <- new_chan
}
