package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

var (
	sshHost     *string
	sshUser     *string
	sshKeyPath  *string
	sshPort     *int
	sshPassword *string
	sshCommand  *string
	sshType     *string
)

func main() {
	sshHost = flag.String("h", "", "(required) EC2 hostname")
	sshUser = flag.String("u", "ec2-user", "(optional) username")
	sshKeyPath = flag.String("i", "/root/admin.pem", "(required) absolute path to aws key file") //ssh id_rsa.id 路径"
	sshPassword = flag.String("P", "", "(optional) sshPassword")
	sshCommand = flag.String("e", "ls", "(required) Command")
	sshType = flag.String("k", "key", "(optional) key|password") //password 或者 key
	sshPort = flag.Int("p", 22, "(optional) 22|22022")
	flag.Parse()
	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            *sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	if *sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(*sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(*sshKeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", *sshHost, *sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	//执行远程命令

	combo, err := session.CombinedOutput(*sshCommand)
	if err != nil {
		//log.Fatal("远程执行cmd 失败", err)
		if v, ok := err.(*ssh.ExitError); ok {
			fmt.Println(v.Msg())
		}
	}
	log.Println("命令输出:")
	fmt.Println(string(combo))

}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	//如果路径以“〜”为前缀，则Expand扩展路径以包括主目录。如果没有以〜为前缀，则按原样返回路径。
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
