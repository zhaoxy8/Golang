package deployexec

import (
	"fmt"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// logger 日志变量
var logger = log.New(os.Stdout, "[k8s]", log.Lshortfile|log.Ldate|log.Ltime)
// 获取执行dp结果的slice切片
var deploymentSlice []*Deployment

// DeploymentConfig 从index.html获取的配置信息
type KubeConfig struct {
	KubeConfig string
	ClientSet *kubernetes.Clientset
}

// NewKubeConfig 构造方法
func NewKubeConfig(kubeconfigform string) *KubeConfig{
	KubeConfigPath := "config/"+ kubeconfigform
	kubeconfig := &KubeConfig{
		KubeConfig: KubeConfigPath,
	}
	//初始化deploymentsClient 字段
	kubeconfig.clinetConfig()
	return kubeconfig
}
// ListNameSpace 获取namespace方法
func (kc *KubeConfig)ListNameSpace() []apiv1.Namespace{
	//logger.Printf("Listing Namespaces in k8s:\n")
	list, err := kc.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	//for _, d := range list.Items {
	//	logger.Printf(" * %s (%s )\n", d.Name, d.Status.Phase)
	//}
	return list.Items
}

func (kc *KubeConfig)clinetConfig() {
	config, err := clientcmd.BuildConfigFromFlags("", kc.KubeConfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	kc.ClientSet = clientset
}

// DeploymentConfig 从index.html获取的配置信息
type DeploymentConfig struct {
	KubeConfig *KubeConfig
	Image string
	Deployment string
	Replicas int32
	Namespace string
	Containersname string
	Containerport int32
	Containerportname string
	Command []string
	Protocol apiv1.Protocol
	ImagePullPolicy apiv1.PullPolicy
	Restartpolicy string
	ImagePullSecrets string
	LabelsMap map[string]string
	Lablekey1 string
	Lablevalue1 string
	DeploymentsClient v1.DeploymentInterface
}
// NewDeploymentConfig 构造方法 传递的参数不全，需要补充
func NewDeploymentConfig(kubeConfig *KubeConfig,Image string,Command []string,Deployment string,Replicas int32,Namespace string) *DeploymentConfig{
	deploymentConfig := &DeploymentConfig{
		KubeConfig: kubeConfig,
		Image: Image,
		Command: Command,
		Deployment: Deployment,
		Replicas: Replicas,
		Namespace: Namespace,
	}
	//初始化deploymentsClient 字段
	deploymentConfig.clinetConfig()
	return deploymentConfig
}

// Deployment 详细信息
type Deployment struct {
	Image string
	Command []string
	Name string
	Replicas int32
	Namespace string
	Status int32
	Labels map[string]string
}
func NewDeployment(Image string,Command []string,Name string,Replicas int32,Namespace string,Status int32,Labels map[string]string) *Deployment{
	deployment := &Deployment{
		Image: Image,
		Command: Command,
		Name: Name,
		Replicas: Replicas,
		Namespace: Namespace,
		Status:Status,
		Labels: Labels,
	}
	return deployment
}


func (dc *DeploymentConfig)clinetConfig() {
	deploymentsClient := dc.KubeConfig.ClientSet.AppsV1().Deployments(dc.Namespace)
	dc.DeploymentsClient = deploymentsClient
}
func (dc *DeploymentConfig)CreateDeployment()(message string){
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dc.Deployment,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: dc.int32Ptr(dc.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: dc.LabelsMap,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: dc.LabelsMap,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dc.Containersname,
							Image: dc.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          dc.Containerportname,
									Protocol:      dc.Protocol,
									ContainerPort: dc.Containerport,
								},
							},
							ImagePullPolicy: dc.ImagePullPolicy,
							Command: dc.Command,
						},
					},
				},
			},
		},
	}
	// Create Deployment
	logger.Println("Creating deployment...")
	result, err := dc.DeploymentsClient.Create(deployment)
	if err != nil {
		return fmt.Sprintf(err.Error())
	}
	return  fmt.Sprintf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}
func (dc *DeploymentConfig)ListDeployment() []*Deployment{
	logger.Printf("Listing deployments in namespace %q:\n", dc.Namespace)
	list, err := dc.DeploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	//初始化一个deployment的切片
	deploymentSlice = make([]*Deployment,0)
	//循环每个deployment构造所需参数的deployment结构体，
	//目前只支持1个deployment中包含1个containers
	for _, d := range list.Items {
		deployment := NewDeployment(
			d.Spec.Template.Spec.Containers[0].Image,
			d.Spec.Template.Spec.Containers[0].Command,
			d.Name,
			*d.Spec.Replicas,
			d.Namespace,
			d.Status.ReadyReplicas,
			d.Spec.Selector.MatchLabels)
		//把deployment结构体放到切片中
		deploymentSlice = append(deploymentSlice, deployment)
		logger.Printf(" %d (%s replicas) %s\n",deployment.Status, deployment.Namespace,deployment.Image)
	}
	//返回deployment切片
	return deploymentSlice
}
func (dc *DeploymentConfig)UpdateDeployment() (result string){
	//logger.Println("Updating deployment...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := dc.DeploymentsClient.Get(dc.Deployment, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}

		result.Spec.Replicas = dc.int32Ptr(dc.Replicas)          // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = dc.Image // change nginx version
		_, updateErr := dc.DeploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		return fmt.Sprintf("Update failed: %v", retryErr)
	}
	//logger.Println("Updated deployment...")
	return fmt.Sprintf("Updated deployment %s successed",dc.Deployment)
}
func (dc *DeploymentConfig)DeleteDeployment() (result string){
	//logger.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := dc.DeploymentsClient.Delete(dc.Deployment, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Sprintf("Delete failed: %v",err.Error())
	}
	//logger.Println("Deleted deployment.")
	return fmt.Sprintf("Deleted deployment %s successed",dc.Deployment)
}
func (dc *DeploymentConfig)int32Ptr(i int32) *int32 { return &i }
func (dc *DeploymentConfig)Run(){
	//dc.ListDeployment()
	dc.CreateDeployment()
}

func ExecComm(c *gin.Context){
	kubeconfigform := c.PostForm("kubeconfig")
	kubeConfig := NewKubeConfig(kubeconfigform)
	namespace := c.PostForm("namespace")
	deployment := c.PostForm("deployment")
	command := strings.Split(c.PostForm("command")," ")
	image := c.PostForm("image")
	replicas,_ := strconv.Atoi(c.PostForm("replicas"))
	//MTInstance := c.PostForm("selectInstance")
	deploymentConfig := NewDeploymentConfig(kubeConfig,image,command,deployment, int32(replicas),namespace)
	fmt.Println(deploymentConfig)
	logger.Println(deploymentConfig)
	deploymentConfig.Run()
	c.HTML(http.StatusOK,"posts/base.html",gin.H{
		"HostResultSlice":1111,
	})
}


func ListNameSpace(c *gin.Context){
	kubeconfigform := c.PostForm("kubeconfig")
	//MTInstance := c.PostForm("selectInstance")
	kubeConfig := NewKubeConfig(kubeconfigform)
	namespaces := kubeConfig.ListNameSpace()
	c.HTML(http.StatusOK,"system/namespace.html",gin.H{
		"namespaces":namespaces,
	})
}

func ListDeployment(c *gin.Context){
	kubeconfigform := c.PostForm("kubeconfig")
	kubeConfig := NewKubeConfig(kubeconfigform)
	namespace := c.PostForm("namespace")
	deploymentConfig := &DeploymentConfig{
		KubeConfig: kubeConfig,
		Namespace: namespace,
	}
	deploymentConfig.clinetConfig()
	deploymentSlice := deploymentConfig.ListDeployment()
	c.HTML(http.StatusOK,"system/deployment-list.html",gin.H{
		"deploymentSlice":deploymentSlice,
	})
}

func DeleteDeployment(c *gin.Context){
	kubeconfigform := c.PostForm("kubeconfig")
	kubeConfig := NewKubeConfig(kubeconfigform)
	namespace := c.PostForm("namespace")
	deployment := c.PostForm("deployment")
	deploymentConfig := &DeploymentConfig{
		KubeConfig: kubeConfig,
		Namespace: namespace,
		Deployment: deployment,
	}
	deploymentConfig.clinetConfig()
	result := deploymentConfig.DeleteDeployment()
	c.HTML(http.StatusOK,"system/deployment-delete.html",gin.H{
		"result":result,
	})
}

func UpdateDeployment(c *gin.Context){
	kubeconfigform := c.PostForm("kubeconfig")
	kubeConfig := NewKubeConfig(kubeconfigform)
	namespace := c.PostForm("namespace")
	deployment := c.PostForm("deployment")
	image := c.PostForm("image")
	replicas,_ := strconv.Atoi(c.PostForm("replicas"))
	deploymentConfig := &DeploymentConfig{
		KubeConfig: kubeConfig,
		Namespace: namespace,
		Deployment: deployment,
		Replicas: int32(replicas),
		Image: image,
	}
	deploymentConfig.clinetConfig()
	result := deploymentConfig.UpdateDeployment()
	c.HTML(http.StatusOK,"system/deployment-update.html",gin.H{
		"result":result,
	})
}

func CreateDeployment(c *gin.Context){
	kubeconfigform := c.PostForm("kubeconfig")
	kubeConfig := NewKubeConfig(kubeconfigform)
	namespace := c.PostForm("namespace")
	deployment := c.PostForm("deployment")
	image := c.PostForm("image")
	replicas,_ := strconv.Atoi(c.PostForm("replicas"))
	containername := c.PostForm("containername")
	containerport,_ := strconv.Atoi(c.PostForm("containerport"))
	protocol := apiv1.Protocol(c.PostForm("protocol"))
	containerportname := c.PostForm("containerportname")
	//command默认是空的Slice，如果有shell命令需要执行才需要构造新的命令行Slice
	command := []string{}
	if c.PostForm("command") != "" {
		command = []string{"sh", "-c"}
		command = append(command, c.PostForm("command"))
	}
	//logger.Println(command)
	imagePullPolicy := apiv1.PullPolicy(c.PostForm("imagePullPolicy"))
	restartpolicy := c.PostForm("restartpolicy")
	imagePullSecrets := c.PostForm("imagePullSecrets")
	lablekey1 := c.PostForm("lablekey1")
	lablevalue1 := c.PostForm("lablevalue1")
	lablekey2 := c.PostForm("lablekey2")
	lablevalue2 := c.PostForm("lablevalue2")
	lablekey3 := c.PostForm("lablekey3")
	lablevalue3 := c.PostForm("lablevalue3")

	deploymentConfig := &DeploymentConfig{
		KubeConfig: kubeConfig,
		Namespace: namespace,
		Deployment: deployment,
		Replicas: int32(replicas),
		Image: image,
		Containersname: containername,
		Containerport: int32(containerport),
		Containerportname: containerportname,
		Command: command,
		Protocol: protocol,
		ImagePullPolicy: imagePullPolicy,
		Restartpolicy: restartpolicy,
		ImagePullSecrets: imagePullSecrets,
		LabelsMap: make(map[string]string),
		Lablekey1: lablekey1,
		Lablevalue1: lablevalue1,
	}
	//必须添加lable
	deploymentConfig.LabelsMap[lablekey1] = lablevalue1
	//如果lable2,lable3不为空，也添加到deployment的lable中
	//for i := 2; i < 3; i++ {
	//	key :=  fmt.Sprintf("lablekey%d",i)
	//	value := fmt.Sprintf("lablevalue%d",i)
	//	if key != "" && value != "" {
	//		deploymentConfig.LabelsMap[key] = value
	//	}
	//}
	if lablekey2 != "" && lablevalue2 != "" {
		deploymentConfig.LabelsMap[lablekey2] = lablevalue2
	}
	if lablekey3 != "" && lablevalue3 != "" {
		deploymentConfig.LabelsMap[lablekey3] = lablevalue3
	}
	logger.Println(deploymentConfig.LabelsMap)
	deploymentConfig.clinetConfig()
	result := deploymentConfig.CreateDeployment()
	c.HTML(http.StatusOK,"system/deployment-create.html",gin.H{
		"result":result,
	})
}