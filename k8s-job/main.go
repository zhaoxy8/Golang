package main

import (

	"os"
	"fmt"
	"flag"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)


func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != ""{
		// 配置 k8s 集群 kubeconfig 配置文件，默认位置 $HOME/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config,err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connect kubernetes Successful！！！")
}
