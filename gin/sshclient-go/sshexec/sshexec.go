package sshexec

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)


// logger 日志变量
var logger = log.New(os.Stdout, "[SSH]", log.Lshortfile|log.Ldate|log.Ltime)
// HostResultSlice 切片
var HostResultSlice []*HostResult
// HostResult ssh命令执行之后的结果结构体
type HostResult struct {
	Host string
	SshResult string
}
// HostResult 构造方法
func NewHostResult(ip string,result string) *HostResult{
	return &HostResult{
		Host:ip,
		SshResult: result,
	}
}
// HostConfig 从index.html获取的配置信息
type HostConfig struct {
	IpList []string
	MtInstance string
	Command string
	Username string
	Port int
	Key string
	Number int
	SshConfig *ssh.ClientConfig
	ComboChan chan *HostResult
	Wg sync.WaitGroup
}
// NewHostConfig 构造方法
func NewHostConfig(IpList []string,MtInstance string,Command string,Username string,Port int,Key string,Number int) *HostConfig{
	hostConfig := &HostConfig{
		IpList: IpList,
		MtInstance: MtInstance,
		Command: Command,
		Username: Username,
		Port: Port,
		Key: Key,
		Number: Number,
		ComboChan: make(chan *HostResult,100),//可以缓存100个ssh结果
	}
	//初始化SshConfig 字段
	hostConfig.sshConfig()
	return hostConfig
}

func (hc *HostConfig)publicKeyAuthFunc() ssh.AuthMethod{
	//如果路径以“〜”为前缀，则Expand扩展路径以包括主目录。如果没有以〜为前缀，则按原样返回路径。
	keyPath, err := homedir.Expand(hc.Key)
	if err != nil {
		logger.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logger.Fatal("ssh key file read failed", err)
		return nil //
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func (hc *HostConfig)sshConfig(){
	//创建ssh登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            hc.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	//每台主机都使用sshkey登录
	config.Auth = []ssh.AuthMethod{hc.publicKeyAuthFunc()}
	//给结构体的SshConfig赋值
	hc.SshConfig = config
}

func (hc *HostConfig)execCom(ip string){
	defer hc.Wg.Done()
	addr := fmt.Sprintf("%s:%d", ip, hc.Port)
	sshClient, err := ssh.Dial("tcp", addr, hc.SshConfig)
	if err != nil {
		logger.Println("创建ssh client 失败", err)
		return
	}
	defer sshClient.Close()
	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		logger.Fatal("创建ssh session 失败", err)
		return
	}
	defer session.Close()
	//执行远程命令
	combo, err := session.CombinedOutput(hc.Command)
	if err != nil {
		logger.Println("远程执行cmd 失败", err)
		if v, ok := err.(*ssh.ExitError); ok {
			logger.Println(v.Lang())
			hostResult := NewHostResult(ip,v.String())
			//把错误结果也存储到channel中
			hc.ComboChan <- hostResult
			return
		}
	}
	//logger.Println("命令输出:")
	//logger.Println(string(combo))
	//构造此IP的执行结果构造成结构体
	hostResult := NewHostResult(ip,string(combo))
	//把结果存储到channel中
	hc.ComboChan <- hostResult
	return
}

func (hc *HostConfig)run(){
	HostResultSlice = make([]*HostResult,0)
	//dial 获取ssh client
	listlen := len(hc.IpList)
	hc.Wg.Add(listlen)
	for _,ip := range hc.IpList{
		go hc.execCom(ip)
	}
	//等待所有goroutine执行完成
	hc.Wg.Wait()
	//所有goroutine执行完成后关闭通道,如果不关闭通道，range是阻塞状态,网页一直转圈
	close(hc.ComboChan)
	//从channel中获取数据，100个缓存，如果ComboChan为空就退出循环，因为通道已经关闭
	for result := range hc.ComboChan{
		logger.Printf("host: %s,result: %s\n",result.Host,result.SshResult)
		//循环从ComboChan中取数据加入到切片中
		HostResultSlice = append(HostResultSlice,result)
	}
	//
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
	c.HTML(http.StatusOK,"posts/base.html",gin.H{
		"HostResultSlice":HostResultSlice,
	})
}
