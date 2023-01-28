package contract

const AppKey = "ext:app"

type App interface {
	//表示当前app的唯一ID
	AppId() string
	//当前版本
	Version() string
	//项目基础目录
	BaseFolder() string
	//配置文件目录
	ConfigFolder() string
	//日志文件目录
	LogFolder() string
	//服务提供者目录
	ProviderFolder() string
	//中间件目录
	MiddlewareFolder() string
	//命令目录
	CommandFolder() string
	//运行时信息
	RuntimeFolder() string
	//测试信息
	TestFolder() string
	//存储目录
	StorageFolder() string

	//AppFolder 定义业务代码所在的目录，用于监控文件变更使用
	AppFolder() string
	//LoadAppConfig 加载新的AppConfig，key为对应的函数转为小写下划线，比如ConfigFolder => config_folder
	LoadAppConfig(configs map[string]string)
}
