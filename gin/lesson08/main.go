package main

import (
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
	route.Static("/static", "./static")
	//模板解析
	route.LoadHTMLGlob("templates/**/*") //**代表目录
	//模板渲染
	route.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK,"posts/index.tmpl",gin.H{
			"title":"<a href='https://liwenzhou.com'>李文周的博客</a>",
		})
	})
	//模板渲染
	route.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK,"users/index.tmpl",gin.H{
			"title":"users/index",
		})
	})
	route.Run(":9090")
}
