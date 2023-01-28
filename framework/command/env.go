package command

import (
	"fmt"
	"github.com/lackone/gin-ext/framework/cobra"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/util"
)

func initEnvCmd() *cobra.Command {
	EnvCmd.AddCommand(EnvListCmd)
	return EnvCmd
}

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "获取当前的App环境",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		env := container.MustMake(contract.EnvKey).(contract.Env)
		fmt.Println("APP_ENV:", env.AppEnv())
		return nil
	},
}

var EnvListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取所有环境变量",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		env := container.MustMake(contract.EnvKey).(contract.Env)
		envs := env.All()
		outs := [][]string{}
		for k, v := range envs {
			outs = append(outs, []string{k, v})
		}
		util.PrettyPrint(outs)
		return nil
	},
}
