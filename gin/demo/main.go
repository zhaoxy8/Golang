package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建一个默认的路由引擎
	// 默认使用了2个中间件Logger(),Recovery()
	r := gin.Default()
	// GET：请求方式；/hello：请求的路径
	// 一、当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("/hello", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})

	})
	// 二、RESful
	r.POST("/Post", getting)
	// 三、API参数
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK,name)
	})
	// 四 URL参数


	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run(":8000")
}

func getting(c *gin.Context)  {
	c.String(http.StatusOK,"Post")
}
