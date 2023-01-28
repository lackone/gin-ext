package env

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/lackone/gin-ext/framework/contract"
	"io"
	"os"
	"path"
	"strings"
)

type ExtEnv struct {
	folder string            //.env所在目录
	envs   map[string]string //所有环境变量
}

func NewExtEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("param error")
	}

	folder := params[0].(string)

	env := &ExtEnv{
		folder: folder,
		envs:   map[string]string{"APP_ENV": contract.EnvDev},
	}

	//解析folder/.env文件
	file := path.Join(folder, ".env")
	fp, err := os.Open(file)
	if err == nil {
		defer fp.Close()
		reader := bufio.NewReader(fp)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			str := bytes.SplitN(line, []byte{'='}, 2)
			if len(str) < 2 {
				continue
			}
			env.envs[string(str[0])] = string(str[1])
		}
	}

	// 获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, val := range os.Environ() {
		e := strings.SplitN(val, "=", 2)
		if len(e) < 2 {
			continue
		}
		env.envs[e[0]] = e[1]
	}

	return env, nil
}

func (e *ExtEnv) AppEnv() string {
	return e.Get("APP_ENV")
}

func (e *ExtEnv) IsExist(key string) bool {
	_, ok := e.envs[key]
	return ok
}

func (e *ExtEnv) Get(key string) string {
	if v, ok := e.envs[key]; ok {
		return v
	}
	return ""
}

func (e *ExtEnv) All() map[string]string {
	return e.envs
}
