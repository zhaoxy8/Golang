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
	route.GET("/profile.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/profile.html",nil)
	})
	//查询namespace执行器
	route.GET("/namespace.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/namespace.html",nil)
	})
	route.POST("/namespacesearch",deployexec.ListNameSpace )
	//分组路由(Grouping Routes) 查询deployment执行器
	//dp := route.Group("/deployment")
	//{
	//	dp.GET("/list", func(c *gin.Context) {
	//		c.HTML(http.StatusOK,"system/deployment-list.html",nil)
	//	})
	//	dp.GET("/create", func(c *gin.Context) {
	//		c.HTML(http.StatusOK,"system/create-deployment-list.html",nil)
	//	})
	//	dp.GET("/update", func(c *gin.Context) {
	//		c.HTML(http.StatusOK,"system/update-deployment-list.html",nil)
	//	})
	//	dp.GET("/delete", func(c *gin.Context) {
	//		c.HTML(http.StatusOK,"system/delete-deployment-list.html",nil)
	//	})
	//}
	route.GET("/deployment-list.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/deployment-list.html",nil)
	})
	route.POST("/deployment-list",deployexec.ListDeployment)

	route.GET("/deployment-create.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/deployment-create.html",nil)
	})
	route.POST("/deployment-create",deployexec.CreateDeployment)
	route.GET("/deployment-update.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/deployment-update.html",nil)
	})
	route.POST("/deployment-update",deployexec.UpdateDeployment)

	route.GET("/deployment-delete.html", func(c *gin.Context) {
		c.HTML(http.StatusOK,"system/deployment-delete.html",nil)
	})
	route.POST("/deployment-delete",deployexec.DeleteDeployment)
	route.Run("0.0.0.0:9090")
}

