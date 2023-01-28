package command

import (
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/lackone/gin-ext/framework/cobra"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/util"
	"github.com/sevlyar/go-daemon"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var cronDaemon = false

func initCronCmd() *cobra.Command {
	CronStartCmd.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start serve daemon")

	CronCmd.AddCommand(CronListCmd)
	CronCmd.AddCommand(CronStartCmd)
	CronCmd.AddCommand(CronStopCmd)
	CronCmd.AddCommand(CronStateCmd)
	CronCmd.AddCommand(CronRestartCmd)

	return CronCmd
}

var CronCmd = &cobra.Command{
	Use:   "cron",
	Short: "cron 定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var CronListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		specs := cmd.Root().CronSpecs
		ps := [][]string{}
		for _, spec := range specs {
			line := []string{spec.Type, spec.Spec, spec.Cmd.Use, spec.Cmd.Short, spec.ServiceName}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

var CronStartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		app := container.MustMake(contract.AppKey).(contract.App)
		//设置cron日志地址和进程ID地址
		pidFile := filepath.Join(app.RuntimeFolder(), "cron.pid")
		logFile := filepath.Join(app.LogFolder(), "cron.log")
		curFolder := app.BaseFolder()

		if cronDaemon {
			cronCtx := &daemon.Context{
				// 设置pid文件
				PidFileName: pidFile,
				PidFilePerm: 0664,
				// 设置日志文件
				LogFileName: logFile,
				LogFilePerm: 0640,
				// 设置工作路径
				WorkDir: curFolder,
				// 设置所有设置文件的mask，默认为750
				Umask: 027,
				// 子进程的参数，按照这个参数设置，子进程的命令为 ./go-ext cron start --daemon=true
				Args: []string{"", "cron", "start", "--daemon=true"},
			}
			child, err := cronCtx.Reborn()
			if err != nil {
				return err
			}
			if child != nil {
				// 父进程直接打印启动成功信息，不做任何操作
				fmt.Println("cron serve started, pid:", child.Pid)
				fmt.Println("log file:", logFile)
				return nil
			}
			// 子进程执行Cron.Run
			defer cronCtx.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("go-ext cron")
			cmd.Root().Cron.Run()
			return nil
		}

		// 没有守护进程模式
		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := os.WriteFile(pidFile, []byte(content), 0664)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("go-ext cron")
		cmd.Root().Cron.Run()
		return nil
	},
}

var CronStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		app := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		pidFile := filepath.Join(app.RuntimeFolder(), "cron.pid")

		pid, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}

		if pid != nil && len(pid) > 0 {
			pid, err := strconv.Atoi(string(pid))
			if err != nil {
				return err
			}
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := os.WriteFile(pidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}

var CronStateCmd = &cobra.Command{
	Use:   "state",
	Short: "cron常驻进程状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		app := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		pidFile := filepath.Join(app.RuntimeFolder(), "cron.pid")

		pid, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}

		if pid != nil && len(pid) > 0 {
			pid, err := strconv.Atoi(string(pid))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("cron server started, pid:", pid)
				return nil
			}
		}
		fmt.Println("no cron server start")
		return nil
	},
}

var CronRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "重启cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		app := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		pidFile := filepath.Join(app.RuntimeFolder(), "cron.pid")

		pid, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}

		if pid != nil && len(pid) > 0 {
			pid, err := strconv.Atoi(string(pid))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}
				// check process closed
				for i := 0; i < 10; i++ {
					if util.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("kill process:" + strconv.Itoa(pid))
			}
		}

		cronDaemon = true
		return CronStartCmd.RunE(cmd, args)
	},
}
