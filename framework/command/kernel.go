package command

import "github.com/lackone/gin-ext/framework/cobra"

func AddKernelCommands(rootCmd *cobra.Command) {
	//挂载app命令
	rootCmd.AddCommand(initAppCmd())
	//挂载cron命令
	rootCmd.AddCommand(initCronCmd())
	//挂载env命令
	rootCmd.AddCommand(initEnvCmd())
	//挂载config命令
	rootCmd.AddCommand(initConfigCmd())
	//挂载build命令
	rootCmd.AddCommand(initBuildCmd())
}
