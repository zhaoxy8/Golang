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
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"time"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func redeployDeployment(deployment ,namespace string)  {

		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("ecommerce").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in ecommerce the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		//_, err = clientset.CoreV1().Pods("ecommerce").Get(context.TODO(), "ecommerce-consumer-campaign", metav1.GetOptions{})
		_, err = clientset.AppsV1().Deployments("ecommerce").Get(context.TODO(), deployment, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod ecommerce-consumer-campaign not found in ecommerce namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found  pod %s in ecommerce namespace\n",deployment)
		}
		//namespace := "ecommerce"
		//2021-05-14T11:05:32
		//now := time.Now()
		//nowstr := now.Format("2006-01-02T15:04:05.000")
		deploymentsClient := clientset.AppsV1().Deployments(namespace)

		//deployment := &appsv1.Deployment{
		//	ObjectMeta: metav1.ObjectMeta{
		//		Name: "demo-deployment",
		//	},
		//	Spec: appsv1.DeploymentSpec{
		//		Replicas: int32Ptr(2),
		//		Selector: &metav1.LabelSelector{
		//			MatchLabels: map[string]string{
		//				"app": "demo",
		//			},
		//		},
		//		Template: apiv1.PodTemplateSpec{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Labels: map[string]string{
		//					"app": "demo",
		//				},
		//				Annotations: map[string]string{
		//					"cattle.io/timestamp":nowstr,
		//				},
		//			},
		//			Spec: apiv1.PodSpec{
		//				Containers: []apiv1.Container{
		//					{
		//						Name:  "web",
		//						Image: "nginx:1.12",
		//						Ports: []apiv1.ContainerPort{
		//							{
		//								Name:          "http",
		//								Protocol:      apiv1.ProtocolTCP,
		//								ContainerPort: 80,
		//							},
		//						},
		//					},
		//				},
		//			},
		//		},
		//	},
		//}
		//
		//// Create Deployment
		//fmt.Println("Creating deployment...")
		//result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		//if err != nil {
		//	panic(err)
		//}
		//
		//fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
		//fmt.Printf("Created deployment Annotations %q.\n", result.GetObjectMeta().GetLabels())
		//
		//// Update Deployment
		//prompt()
		fmt.Println("Updating deployment...")
		//    You have two options to Update() this Deployment:
		//
		//    1. Modify the "deployment" variable and call: Update(deployment).
		//       This works like the "kubectl replace" command and it overwrites/loses changes
		//       made by other clients between you Create() and Update() the object.
		//    2. Modify the "result" returned by Get() and retry Update(result) until
		//       you no longer get a conflict error. This way, you can preserve changes made
		//       by other clients between Create() and Update(). This is implemented below
		//			 using the retry utility package included with client-go. (RECOMMENDED)
		//
		// More Info:
		// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
		type ann map[string]string
		annotations := make(ann,1)

		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			result, getErr := deploymentsClient.Get(context.TODO(), deployment, metav1.GetOptions{})
			if getErr != nil {
				panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
			}

			//result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
			oldannotations := result.Spec.Template.Annotations
			fmt.Printf("oldannotations%s",oldannotations)
			now := time.Now()
			nowstr := now.Format("2006-01-02T15:04:05.000")
			annotations["cattle.io/timestamp"]=nowstr
			result.Spec.Template.Annotations = annotations                 // change annotations timestamp
			//result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
			_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			panic(fmt.Errorf("Update failed: %v", retryErr))
		}
		fmt.Println("Updated deployment...")

		// List Deployments
		//prompt()
		fmt.Printf("Listing deployments in namespace %q:\n", namespace)
		//list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
		dep, err := deploymentsClient.Get(context.TODO(), deployment, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf(" * %s (%d replicas)\n", (*dep).Name, *dep.Spec.Replicas)
		// Delete Deployment
		//prompt()
		//fmt.Println("Deleting deployment...")
		//deletePolicy := metav1.DeletePropagationForeground
		//if err := deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		//	PropagationPolicy: &deletePolicy,
		//}); err != nil {
		//	panic(err)
		//}
		//fmt.Println("Deleted deployment.")

}

func main(){
	for{
		go redeployDeployment("ecommerce-consumer-task","ecommerce")
		time.Sleep(180 * time.Second)
	}
}

//func prompt() {
//	fmt.Printf("-> Press Return key to continue.")
//	scanner := bufio.NewScanner(os.Stdin)
//	for scanner.Scan() {
//		break
//	}
//	if err := scanner.Err(); err != nil {
//		panic(err)
//	}
//	fmt.Println()
//}
//
//func int32Ptr(i int32) *int32 { return &i }
