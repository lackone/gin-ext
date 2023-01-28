package services

import (
	"fmt"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/util"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"time"
)

type ExtRotateLog struct {
	ExtLog

	// 日志文件存储目录
	folder string
	// 日志文件名
	file string
}

func NewExtRotateLog(params ...interface{}) (interface{}, error) {
	// 参数解析
	container := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	app := container.MustMake(contract.AppKey).(contract.App)
	conf := container.MustMake(contract.ConfigKey).(contract.Config)

	// 从配置文件中获取folder信息，否则使用默认的LogFolder文件夹
	folder := app.LogFolder()
	if conf.IsExist("log.folder") {
		folder = conf.GetString("log.folder")
	}
	// 如果folder不存在，则创建
	if !util.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	// 从配置文件中获取file信息，否则使用默认的ext.log
	file := "ext.log"
	if conf.IsExist("log.file") {
		file = conf.GetString("log.file")
	}

	// 从配置文件获取date_format信息
	dateFormat := "%Y%m%d%H"
	if conf.IsExist("log.date_format") {
		dateFormat = conf.GetString("log.date_format")
	}

	linkName := rotatelogs.WithLinkName(filepath.Join(folder, file))
	options := []rotatelogs.Option{linkName}

	// 从配置文件获取rotate_count信息
	if conf.IsExist("log.rotate_count") {
		rotateCount := conf.GetInt("log.rotate_count")
		options = append(options, rotatelogs.WithRotationCount(uint(rotateCount)))
	}

	// 从配置文件获取rotate_size信息
	if conf.IsExist("log.rotate_size") {
		rotateSize := conf.GetInt("log.rotate_size")
		options = append(options, rotatelogs.WithRotationSize(int64(rotateSize)))
	}

	// 从配置文件获取max_age信息
	if conf.IsExist("log.max_age") {
		if maxAgeParse, err := time.ParseDuration(conf.GetString("log.max_age")); err == nil {
			options = append(options, rotatelogs.WithMaxAge(maxAgeParse))
		}
	}

	// 从配置文件获取rotate_time信息
	if conf.IsExist("log.rotate_time") {
		if rotateTimeParse, err := time.ParseDuration(conf.GetString("log.rotate_time")); err == nil {
			options = append(options, rotatelogs.WithRotationTime(rotateTimeParse))
		}
	}

	// 设置基础信息
	log := &ExtRotateLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = folder
	log.file = file

	w, err := rotatelogs.New(fmt.Sprintf("%s.%s", filepath.Join(log.folder, log.file), dateFormat), options...)
	if err != nil {
		return nil, errors.New("new rotatelogs error")
	}
	log.SetOutput(w)
	log.container = container
	return log, nil
}
