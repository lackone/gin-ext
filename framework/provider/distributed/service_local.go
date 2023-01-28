package distributed

import (
	"errors"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type ExtDistributedLocal struct {
	container framework.Container
}

func NewExtDistributedLocal(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("param error")
	}

	container := params[0].(framework.Container)
	return &ExtDistributedLocal{container: container}, nil
}

func (e ExtDistributedLocal) Select(serviceName string, appID string, holdTime time.Duration) (selectAppID string, err error) {
	app := e.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := app.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder, "distributed_"+serviceName)

	// 打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// 尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	// 抢不到文件锁
	if err != nil {
		// 读取被选择的appid
		selectAppIDByt, err := io.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt), err
	}

	// 在一段时间内，选举有效，其他节点在这段时间不能再进行抢占
	go func() {
		defer func() {
			// 释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 释放文件
			lock.Close()
			// 删除文件锁对应的文件
			os.Remove(lockFile)
		}()
		// 创建选举结果有效的计时器
		timer := time.NewTimer(holdTime)
		// 等待计时器结束
		<-timer.C
	}()

	// 这里已经是抢占到了，将抢占到的appID写入文件
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}
	return appID, nil
}
