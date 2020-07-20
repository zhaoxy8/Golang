package deployexec

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"sync"
)

// logger 日志变量
var logger = log.New(os.Stdout, "[SSH]", log.Lshortfile|log.Ldate|log.Ltime)

// DeploymentConfig 从index.html获取的配置信息
type DeploymentConfig struct {
	KubeConfig string
	Image string
	Command string
	Deplyment string
	Replicas int
	Namespace string
	Wg sync.WaitGroup
}

// NewDeploymentConfig 构造方法
func NewDeploymentConfig(KubeConfig string,Image string,Command string,Deplyment string,Replicas int,Namespace string) *DeploymentConfig{
	deploymentConfig := &DeploymentConfig{
		KubeConfig: KubeConfig,
		Image: Image,
		Command: Command,
		Deplyment: Deplyment,
		Replicas: Replicas,
		Namespace: Namespace,
	}
	//初始化SshConfig 字段
	//hostConfig.sshConfig()
	return deploymentConfig
}


func ExecComm(c *gin.Context){
	kubeconfig := c.PostForm("kubeconfig")
	namespace := c.PostForm("namespace")
	deplyment := c.PostForm("deplyment")
	command := c.PostForm("command")
	image := c.PostForm("image")
	replicas,_ := strconv.Atoi(c.PostForm("replicas"))
	//MTInstance := c.PostForm("selectInstance")
	deploymentConfig := NewDeploymentConfig(kubeconfig,image,command,deplyment,replicas,namespace)
	fmt.Println(deploymentConfig)
	logger.Println(deploymentConfig)
	//hostConfig.run()
	//c.HTML(http.StatusOK,"posts/base.html",gin.H{
	//	"HostResultSlice":HostResultSlice,
	//})
}