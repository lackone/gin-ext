package demo

import "github.com/lackone/gin-ext/framework/gin"

func Register(r *gin.Engine) error {
	r.GET("/demo/demo", func(c *gin.Context) {
		c.JSON(200, "test")
	})

	return nil
}
