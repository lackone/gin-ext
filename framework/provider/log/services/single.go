package services

import (
	"errors"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/util"
	"os"
	"path/filepath"
)

type ExtSingleLog struct {
	ExtLog

	folder string
	file   string
}

func NewExtSingleLog(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	app := container.MustMake(contract.AppKey).(contract.App)
	conf := container.MustMake(contract.ConfigKey).(contract.Config)

	log := &ExtSingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	folder := app.LogFolder()
	if conf.IsExist("log.folder") {
		folder = conf.GetString("log.folder")
	}
	log.folder = folder
	if !util.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	log.file = "ext.log"
	if conf.IsExist("log.file") {
		log.file = conf.GetString("log.file")
	}

	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.New("open log file err")
	}

	log.SetOutput(fd)
	log.container = container

	return log, nil
}
