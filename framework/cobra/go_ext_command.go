package cobra

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/robfig/cron/v3"
	"log"
)

// CronSpec 保存Cron命令的信息，用于展示
type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

// 设置容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

func (c *Command) SetParentNull() {
	c.parent = nil
}

func (c *Command) AddCronCommand(spec string, cmd *Command) {
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New(cron.WithSeconds())
		root.CronSpecs = []CronSpec{}
	}
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})

	//创建一个command
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	cronCmd.SetContainer(root.GetContainer())

	//添加调用函数
	root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}
