/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
	"time"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var clientset *kubernetes.Clientset

func getPodLogs(pod v1.Pod) string {
	var lines int64
	lines = 100
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{TailLines: &lines})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println( "error in opening stream")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		fmt.Println("error in copy information from podLogs to buf")
	}
	str := buf.String()
	fmt.Println(str)

	return str
}
func runLogs(kubeconfig *string,namespace *string){
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods(*namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name,"pointsshop-h5") {
				fmt.Printf(" Pod %s (%s Status)\n", pod.Name,pod.Status.ContainerStatuses[0].Name )
				str := getPodLogs(pod)
				fmt.Println(str)
				time.Sleep(10 * time.Second)
			}
		}

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "web"
		pod := "pointsshop-h5-7c8b68f767-vwd6g"
		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	var kubeconfig *string
	var namespace *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	namespace = flag.String("n", "", "(optional) namespace")
	flag.Parse()

	for index , value := range flag.Args() {
		fmt.Println(index, value)
	}
	var cmdPull = &cobra.Command{
		Use:   "logs ",
		Short: "logs containers",
		Run: func(cmd *cobra.Command, args []string)  {
			runLogs(kubeconfig,namespace)
		},
	}

	var rootCmd = &cobra.Command{Use: "kubectl"}
	rootCmd.PersistentFlags().String("kubeconfig","","absolute path to the kubeconfig file")
	rootCmd.AddCommand(cmdPull)
	rootCmd.Execute()

	fmt.Printf("kubeconfig %s \n", *kubeconfig)
	fmt.Printf("namespace %s \n", *namespace)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h //root
	}
	return os.Getenv("USERPROFILE") // windows
}

