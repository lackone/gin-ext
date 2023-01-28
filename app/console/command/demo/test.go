package demo

import (
	"github.com/lackone/gin-ext/framework/cobra"
	"log"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "测试",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("test....")
		return nil
	},
}
