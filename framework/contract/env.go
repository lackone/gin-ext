package contract

const (
	EnvProd = "prod"
	EnvTest = "test"
	EnvDev  = "dev"
	EnvKey  = "ext:env"
)

type Env interface {
	//获取当前环境，分为dev/test/prod
	AppEnv() string
	//判断环境变量是否存在
	IsExist(string) bool
	//获取环境变量
	Get(string) string
	//获取所有环境变量
	All() map[string]string
}
