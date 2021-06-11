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
package logs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
func RunLogs(kubeconfig string,namespace string){
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
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
	}
}

func aunLogs(kubeconfig string,namespace string) {
	fmt.Printf("kubeconfig %s \n", kubeconfig)
	fmt.Printf("namespace %s \n", namespace)
}

