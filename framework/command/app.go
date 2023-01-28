package command

import (
	"context"
	"github.com/lackone/gin-ext/framework/cobra"
	"github.com/lackone/gin-ext/framework/contract"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initAppCmd() *cobra.Command {
	appCmd.AddCommand(appStartCmd)

	return appCmd
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "app 相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var appStartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动一个app服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		//获取容器
		container := cmd.GetContainer()
		//从容器中获取kernel服务实例
		kernel := container.MustMake(contract.KernelKey).(contract.Kernel)
		engine := kernel.HttpEngine()

		server := &http.Server{
			Handler: engine,
			Addr:    ":8080",
		}

		go func() {
			server.ListenAndServe()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit

		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()

		if err := server.Shutdown(timeout); err != nil {
			log.Fatalln(err)
		}

		return nil
	},
}
