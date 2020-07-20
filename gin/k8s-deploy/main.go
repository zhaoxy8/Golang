package main

import (
	"Golong/gin/k8s-deploy/deployexec"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)


func main() {
	route := gin.Default()
	//自定义模板函数
	route.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML{
			return template.HTML(str)
		},
	})
	//静态文件处理,前面的目录 解析为后面的本地目录名
	route.Static("/static", "./statics")
	//模板解析
	route.LoadHTMLGlob("templates/**/*") //**代表目录
	//模板渲染
	route.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"login/index.html",nil)
	})
	route.POST("/command",deployexec.ExecComm)
	route.Run(":9090")
}

