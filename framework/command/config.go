package command

import (
    "fmt"
    "github.com/kr/pretty"
    "github.com/lackone/gin-ext/framework/cobra"
    "github.com/lackone/gin-ext/framework/contract"
)

func initConfigCmd() *cobra.Command {
    ConfigCmd.AddCommand(ConfigGetCmd)

    return ConfigCmd
}

var ConfigCmd = &cobra.Command{
    Use:   "config",
    Short: "获取配置相关信息",
    RunE: func(c *cobra.Command, args []string) error {
        if len(args) == 0 {
            c.Help()
        }
        return nil
    },
}

var ConfigGetCmd = &cobra.Command{
    Use:   "get",
    Short: "获取某个配置信息",
    RunE: func(c *cobra.Command, args []string) error {
        container := c.GetContainer()
        conf := container.MustMake(contract.ConfigKey).(contract.Config)
        if len(args) != 1 {
            fmt.Println("参数错误")
            return nil
        }
        configPath := args[0]
        val := conf.Get(configPath)
        if val == nil {
            fmt.Println("配置路径 ", configPath, " 不存在")
            return nil
        }

        fmt.Printf("%# v\n", pretty.Formatter(val))
        return nil
    },
}
