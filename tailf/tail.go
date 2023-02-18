package tailf

import (
	"github.com/hpcloud/tail"
)

var (
	Tail *tail.Tail
	err  error
)

func Init(filename string) error {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	Tail, err = tail.TailFile(filename, config)
	if err != nil {
		return err
	}
	return nil
}
