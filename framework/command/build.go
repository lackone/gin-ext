package command

import (
	"fmt"
	"github.com/lackone/gin-ext/framework/cobra"
	"log"
	"os/exec"
)

// build相关的命令
func initBuildCmd() *cobra.Command {
	BuildCmd.AddCommand(BuildSelfCmd)
	BuildCmd.AddCommand(BuildBackendCmd)
	BuildCmd.AddCommand(BuildFrontendCmd)
	BuildCmd.AddCommand(BuildAllCmd)
	return BuildCmd
}

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "编译相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var BuildSelfCmd = &cobra.Command{
	Use:   "self",
	Short: "编译自身命令",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("go: 请在Path路径中先安装go")
		}
		cmd := exec.Command(path, "build", "-o", "gin-ext", "./")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("go build error:")
			fmt.Println(string(out))
			fmt.Println("--------------")
			return err
		}
		fmt.Println("编译自身成功")
		return nil
	},
}

var BuildBackendCmd = &cobra.Command{
	Use:   "backend",
	Short: "使用go编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return BuildSelfCmd.RunE(c, args)
	},
}

var BuildFrontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "使用npm编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		// 获取path路径下的npm命令
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Fatalln("请安装npm在你的PATH路径下")
		}

		// 执行npm run build
		cmd := exec.Command(path, "run", "build")
		// 将输出保存在out中
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("=============  前端编译失败 ============")
			fmt.Println(string(out))
			fmt.Println("=============  前端编译失败 ============")
			return err
		}
		// 打印输出
		fmt.Print(string(out))
		fmt.Println("=============  前端编译成功 ============")
		return nil
	},
}

var BuildAllCmd = &cobra.Command{
	Use:   "all",
	Short: "同时编译前端和后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := BuildFrontendCmd.RunE(c, args)
		if err != nil {
			return err
		}
		err = BuildBackendCmd.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}
