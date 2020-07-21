package main

import (
	"Golang/gin/k8s-deploy/deployexec"
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
		c.HTML(http.StatusOK,"system/index.html",nil)
	})
	route.GET("/graphs.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/graphs.html",nil)
	})
	route.GET("/maps.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/maps.html",nil)
	})
	route.GET("/typography.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/typography.html",nil)
	})

	route.GET("/inbox.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/inbox.html",nil)
	})

	route.GET("/gallery.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/gallery.html",nil)
	})

	route.GET("/layout.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/layout.html",nil)
	})
	route.GET("/forms.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/forms.html",nil)
	})
	route.GET("/validation.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/validation.html",nil)
	})

	route.GET("/404.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/404.html",nil)
	})
	route.GET("/faq.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/faq.html",nil)
	})
	route.GET("/blank.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/blank.html",nil)
	})
	route.GET("/signin.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/signin.html",nil)
	})
	route.GET("/signup.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/signup.html",nil)
	})
	//查询namespace执行器  http://11.81.3.6:9090/namespace/search?kubeconfig=dev-dr-config
	route.GET("/namespace.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/namespace.html",nil)
	})
	route.GET("/namespace/search",deployexec.ListNameSpace )
	route.POST("/command",deployexec.ExecComm)
	route.Run("0.0.0.0:9090")
}

