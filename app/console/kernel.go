package console

import (
	"github.com/lackone/gin-ext/app/console/command/demo"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/cobra"
	"github.com/lackone/gin-ext/framework/command"
	"time"
)

func RunCommand(container framework.Container) error {
	rootCmd := &cobra.Command{
		Use:   "gin-ext",
		Short: "gin-ext 命令",
		Long:  "gin-ext 框架提供的命令行工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	//为根命令设置服务容器
	rootCmd.SetContainer(container)
	//绑定框架的命令
	command.AddKernelCommands(rootCmd)
	//绑定业务的命令
	AddAppCommand(rootCmd)

	return rootCmd.Execute()
}

// 业务的命令
func AddAppCommand(rootCmd *cobra.Command) {
	//rootCmd.AddCronCommand("* * * * * *", demo.TestCmd)

	rootCmd.AddDistributedCronCommand("test", "* * * * * *", demo.TestCmd, time.Second*3)
}
