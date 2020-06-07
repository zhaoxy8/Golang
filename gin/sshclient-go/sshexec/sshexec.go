package sshexec

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// HostConfig
type HostConfig struct {
	IpList []string
	MtInstance string
	Command string
	Username string
	Port int
	Key string
	Number int
}

func NewHostConfig(IpList []string,MtInstance string,Command string,Username string,Port int,Key string,Number int) *HostConfig{
	return &HostConfig{
		IpList: IpList,
		MtInstance: MtInstance,
		Command: Command,
		Username: Username,
		Port: Port,
		Key: Key,
		Number: Number,
	}
}

func (hc *HostConfig)run(){

}

//ExecComm 获取POST参数
func ExecComm(c *gin.Context){
	command := c.PostForm("command")
	username := c.PostForm("username")
	port,_ := strconv.Atoi(c.PostForm("port"))
	key := c.PostForm("key")
	number,_ := strconv.Atoi(c.PostForm("number"))
	//主机列表使用windows换行符进行切分
	iplist := strings.Split(c.PostForm("iplist"),"\r\n")
	MTInstance := c.PostForm("selectInstance")
	hostConfig := NewHostConfig(iplist,MTInstance,command,username,port,key,number)
	hostConfig.run()
	c.JSON(http.StatusOK,gin.H{
		"Command":hostConfig.Command,
		"Username":hostConfig.Username,
		"Port":hostConfig.Port,
		"Key":hostConfig.Key,
		"Number":hostConfig.Number,
		"IpList":hostConfig.IpList,
		"MtInstance":hostConfig.MtInstance,
	})
}
