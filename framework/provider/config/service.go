package config

import (
	"bytes"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ExtConfig struct {
	container framework.Container    //容器
	folder    string                 //配置目录
	delim     string                 //路径分隔符，默认为点
	lock      sync.RWMutex           //配置文件读写锁
	envs      map[string]string      //所有环境变量
	confMaps  map[string]interface{} //配置文件结构，key为文件名
	confRaws  map[string][]byte      //配置文件原始信息
}

func (e *ExtConfig) IsExist(key string) bool {
	return e.find(key) != nil
}

func (e *ExtConfig) Get(key string) interface{} {
	return e.find(key)
}

func (e *ExtConfig) GetBool(key string) bool {
	return cast.ToBool(e.find(key))
}

func (e *ExtConfig) GetInt(key string) int {
	return cast.ToInt(e.find(key))
}

func (e *ExtConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(e.find(key))
}

func (e *ExtConfig) GetTime(key string) time.Time {
	return cast.ToTime(e.find(key))
}

func (e *ExtConfig) GetString(key string) string {
	return cast.ToString(e.find(key))
}

func (e *ExtConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(e.find(key))
}

func (e *ExtConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(e.find(key))
}

func (e *ExtConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(e.find(key))
}

func (e *ExtConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(e.find(key))
}

func (e *ExtConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(e.find(key))
}

func (e *ExtConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(e.find(key))
}

// 通过path来获取某个配置项
func (e *ExtConfig) find(key string) interface{} {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return searchMap(e.confMaps, strings.Split(key, e.delim))
}

func (e *ExtConfig) loadConfigFile(folder, file string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		data, err := os.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		//直接针对文本做环境变量的替换
		data = replace(data, e.envs)
		//解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(data, &c); err != nil {
			return err
		}

		e.confMaps[name] = c
		e.confRaws[name] = data

		//读取app.path中的信息，更新app对应的目录
		if name == "app" && e.container.IsBind(contract.AppKey) {
			if v, ok := c["path"]; ok {
				app := e.container.MustMake(contract.AppKey).(contract.App)
				app.LoadAppConfig(cast.ToStringMapString(v))
			}
		}
	}

	return nil
}

// 删除文件的操作
func (e *ExtConfig) removeConfigFile(folder string, file string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(e.confRaws, name)
		delete(e.confMaps, name)
	}
	return nil
}

// replace 表示使用环境变量maps替换context中的env(xxx)的环境变量
func replace(content []byte, envs map[string]string) []byte {
	if envs == nil {
		return content
	}
	// 直接使用ReplaceAll替换。这个性能可能不是最优，但是配置文件加载，频率是比较低的，可以接受
	for key, val := range envs {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}

func NewExtConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envs := params[2].(map[string]string)

	conf := &ExtConfig{
		container: container,
		folder:    envFolder,
		envs:      envs,
		confRaws:  map[string][]byte{},
		confMaps:  map[string]interface{}{},
		delim:     ".",
		lock:      sync.RWMutex{},
	}

	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		// 这里修改成为不返回错误，是让new方法可以通过
		return conf, nil
	}

	//读取目录下每个文件
	files, err := os.ReadDir(envFolder)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for _, file := range files {
		fileName := file.Name()
		err := conf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(envFolder)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watcher.Events:
				{
					//判断事件发生的类型
					// Create 创建
					// Write 写入
					// Remove 删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						conf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						conf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						conf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watcher.Errors:
				{
					log.Println(err)
					return
				}
			}
		}
	}()

	return conf, nil
}

// 查找某个路径的配置项
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	// 判断是否有下个路径
	next, ok := source[path[0]]
	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}

		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface的map，使用cast进行下value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// 如果是map[string]，直接循环调用
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			// 否则的话，返回nil
			return nil
		}
	}
	return nil
}
