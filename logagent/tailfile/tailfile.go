package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var (
	TailObj *tail.Tail
)

func Init(fileName string) (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	TailObj, err = tail.TailFile(fileName, cfg)
	if err != nil {
		logrus.Error("tailfile: create tailObj for path:%s failed, err:%v\n", fileName, err)
		return
	}
	return
}
