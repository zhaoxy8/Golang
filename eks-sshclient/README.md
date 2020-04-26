 - 实现和ssh命令相同的功能
`golang.org/x/crypto/ssh`
`[root@ip-11-81-1-194 eks-sshclient]# ./eks-sshclient -h 
flag needs an argument: -h
Usage of ./eks-sshclient:
  -P string
        (optional) sshPassword
  -e string
        (required) Command (default "ls")
  -h string
        (required) EC2 hostname
  -i string
        (required) absolute path to aws key file (default "/root/admin.pem")
  -k string
        (optional) key|password (default "key")
  -p int
        (optional) 22|22022 (default 22)
  -u string
        (optional) username (default "ec2-user")`
