package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type SimpleRequest struct {
	Text    SimpleRequestContent `json:"text"`
	Msgtype string               `json:"msgtype"`
}

type SimpleRequestContent struct {
	Content string `json:"content"`
}
//
//var UrlAddress = "https://oapi.dingtalk.com/robot/send?access_token=b291d1fc6728b6be13b786ec44809bb36425a33cec022c7782c98d19356f6f87"
//var Secret = "SECf256c8ae405df3f958b613aa41e723b3a3aeaddff182e69547288dda2fa31a66"
//var NacosUrl = "https://m-nacos-stg.bmw-emall.cn"
//设置全局变量 需要ENV 变量传入替换
var NacosUrl string
var UrlAddress string
var Secret string
var ProjectEnv string
//常量
const URI = "/nacos/v1/ns/instance/list?serviceName="

type Host struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Valid      bool   `json:"valid"`
	Healthy    bool   `json:"healthy"`
	Marked     bool   `json:"marked"`
	InstanceID string `json:"instanceId"`
	Metadata   struct {
		PreservedRegisterSource string `json:"preserved.register.source"`
		Version                 string `json:"version"`
	} `json:"metadata"`
	Enabled     bool    `json:"enabled"`
	Weight      float64 `json:"weight"`
	ClusterName string  `json:"clusterName"`
	ServiceName string  `json:"serviceName"`
	Ephemeral   bool    `json:"ephemeral"`
}
type ServiceAll struct {
	Hosts           []Host `json:"hosts"`
	Dom             string `json:"dom"`
	Name            string `json:"name"`
	CacheMillis     int    `json:"cacheMillis"`
	LastRefTime     int64  `json:"lastRefTime"`
	Checksum        string `json:"checksum"`
	UseSpecifiedURL bool   `json:"useSpecifiedURL"`
	Clusters        string `json:"clusters"`
	Env             string `json:"env"`
	Metadata        struct {
	} `json:"metadata"`
}

func selectAllInstances(servicename string) []Host {
	var service ServiceAll
	//Membership
	serviceUrl := NacosUrl + URI + servicename
	//fmt.Println(NacosUrl)
	resp, err := http.Get(serviceUrl)
	if err != nil {
		fmt.Println("http.Get err=", err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
	}
	//fmt.Println("curl:" + string(bytes))
	json.Unmarshal(bytes, &service)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
	}
	//fmt.Printf("hosts:%#v\n", ins)
	//for i := 0; i < len(service.Hosts); i++ {
	//	fmt.Println(service.Hosts[i].Port, service.Hosts[i].ServiceName, service.Hosts[i].IP)
	//}
	return service.Hosts
}
func sendDingding(msg string) {
	// 设置请求头
	requestBody := SimpleRequest{
		Text: SimpleRequestContent{
			Content: msg,
		},
		Msgtype: "text",
	}
	reqBodyBox, _ := json.Marshal(requestBody)
	body := string(reqBodyBox)

	//  构建 签名
	//  把timestamp+"\n"+密钥当做签名字符串，使用HmacSHA256算法计算签名，然后进行Base64 encode，最后再把签名参数再进行urlEncode，得到最终的签名（需要使用UTF-8字符集）。
	timeStampNow := time.Now().UnixNano() / 1000000
	signStr := fmt.Sprintf("%d\n%s", timeStampNow, Secret)

	hash := hmac.New(sha256.New, []byte(Secret))
	hash.Write([]byte(signStr))
	sum := hash.Sum(nil)

	encode := base64.StdEncoding.EncodeToString(sum)
	urlEncode := url.QueryEscape(encode)

	// 构建 请求 url
	UrlAddress = fmt.Sprintf("%s&timestamp=%d&sign=%s", UrlAddress, timeStampNow, urlEncode)

	// 构建 请求体
	request, _ := http.NewRequest("POST", UrlAddress, strings.NewReader(body))

	// 设置库端口
	client := &http.Client{}

	// 请求头添加内容
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	response, _ := client.Do(request)
	//fmt.Println("response: ", response)

	// 关闭 读取 reader
	defer response.Body.Close()

	// 读取内容
	all, _ := ioutil.ReadAll(response.Body)
	fmt.Println("all: ", string(all))
}
func healthchek() {
	defer func() {
		if info := recover(); info != nil {
			fmt.Println("触发了panic", info)
			return
		} else {
			fmt.Println("程序正常退出")
		}
	}()
	healthMap := make(map[string]string, 100)
	//Membership
	var m2service = [...]string{"membership-consumer-bizapi",
		"membership-provider-marketactivity",
		"membership-provider-memberservice",
		"membership-common-auth",
		"membership-provider-integrationcdp",
		"membership-provider-points",
		"membership-provider-pointsmall",
		"membership-provider-integrationdmo",
		"membership-consumer-pointsmall",
		"membership-task",
		"membership-provider-voucher",
		"membership-provider-msg",
		"membership-consumer-cpapi",
		"membership-provider-cust",
		"membership-consumer-integateapi",
	}

	// 方法1：for循环遍历
	for i := 0; i < len(m2service); i++ {
		//fmt.Println(m2service[i])
		instances := selectAllInstances(m2service[i])
		//fmt.Println(strconv.Itoa(len(instances)))
		//1.如果服务可用节点数小于等于1报警
		if len(instances) <= 1 {
			sendDingding(ProjectEnv+" Health Check Error: "+m2service[i]+" Available nodes is [ "+ strconv.Itoa(len(instances)) +" ] ")
		}
		//2.检查服务节点健康状态
		for i := 0; i < len(instances); i++ {
			//每个实例健康检查结果 1 健康 0 不健康
			var healthResult string
			var healthService string
			var bytes []byte
			//Membership
			url := "/publicApi/health"
			if instances[i].ServiceName == "membership-common-auth" {
				url = "/uaa/publicApi/health"
				//url = "/uaa/publicApi/health1"
			}
			//Ecommerce
			//url := "/test/health"
			//fmt.Printf("%s,%s,%d \n", instances[i].ServiceName, instances[i].IP, instances[i].Port)
			//port := strconv.FormatUint(instances[i].Port, 10) uint64类型需要这个转换
			url = "http://" + instances[i].IP + ":" + strconv.Itoa(instances[i].Port) + url
			//fmt.Println(url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("http.Get err=", err)
				//3.连接拒绝的时候报警
				sendDingding(fmt.Sprintf(ProjectEnv+" Health Check Error: ",err))
			}
			//fmt.Println(resp.StatusCode,200)
			if resp.StatusCode != 200 {
				//bytes, err = ioutil.ReadAll(resp.Body)
				//if err != nil {
				//	fmt.Println("ioutil.ReadAll err=", err)
				//}
				//fmt.Println(string(bytes))
				bytes = []byte("error")
			} else {
				var err error
				bytes, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("ioutil.ReadAll err=", err)
				}
			}

			//fmt.Println(string(bytes))
			//container_health{container_name="membership-provider-cust",ip="10.85.46.102",port="10020"} 0
			//container_health{container_name="membership-provider-pointsmall",ip="",ip="10.85.46.103",port="20020"} 1
			healthService = fmt.Sprintf("container_health{container_name=\"%s\",ip=\"%s\",port=\"%d\"}", instances[i].ServiceName, instances[i].IP, instances[i].Port)
			if string(bytes) == "SUCCESS" {
				healthResult = "1"
				//sendDingding(" Health Check Error: "+healthService)
			} else {
				healthResult = "0"
				//发送钉钉报警
				sendDingding(ProjectEnv+" Health Check Error: "+healthService)
			}
			//构造map 给prom使用的数据，已经不需要
			healthMap[healthService] = healthResult
		}
	}
}


func main() {
	//从环境变量中取值赋值
	NacosUrl = os.Getenv("NACOS_CLUSTER")
	UrlAddress = os.Getenv("DINGDING_URLADDRESS")
	Secret = os.Getenv("DINGDING_SECRET")
	ProjectEnv = os.Getenv("PROJECT_ENV")
	var count = 0
	for {
		fmt.Println("sevice-healthcheck start:", count)
		healthchek()
		time.Sleep(time.Second * 120)
		count++
	}
}
