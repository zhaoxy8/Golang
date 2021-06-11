package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

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
	//Ecom
	//serviceUrl := "https://nacos-internal-dev.bmw-emall.cn/nacos/v1/ns/instance/list?serviceName=" + servicename
	//Membership
	serviceUrl := "https://m-nacos-stg.bmw-emall.cn/nacos/v1/ns/instance/list?serviceName=" + servicename
	//fmt.Println(serviceUrl)
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
func healthchek(c *gin.Context) {
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
		//fmt.Println(instances)
		for i := 0; i < len(instances); i++ {
			//每个实例健康检查结果 1 健康 0 不健康
			var healthResult string
			var healthService string
			var bytes []byte
			//Membership
			url := "/publicApi/health"
			if instances[i].ServiceName == "membership-common-auth"{
				url = "/uaa/publicApi/health"
			}
			//Ecommerce
			//url := "/test/health"
			fmt.Printf("%s,%s,%d \n", instances[i].ServiceName, instances[i].IP, instances[i].Port)
			//port := strconv.FormatUint(instances[i].Port, 10) uint64类型需要这个转换
			url = "http://" + instances[i].IP + ":" + strconv.Itoa(instances[i].Port) + url
			fmt.Println(url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("http.Get err=", err)
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

			fmt.Println(string(bytes))
			//container_health{container_name="membership-provider-cust",ip="10.85.46.102",port="10020"} 0
			//container_health{container_name="membership-provider-pointsmall",ip="",ip="10.85.46.103",port="20020"} 1
			healthService = fmt.Sprintf("container_health{container_name=\"%s\",ip=\"%s\",port=\"%d\"}", instances[i].ServiceName, instances[i].IP, instances[i].Port)
			if string(bytes) == "SUCCESS" {
				healthResult = "1"
			} else {
				healthResult = "0"
			}
			healthMap[healthService] = healthResult

		}
	}
	//c.String(http.StatusOK, "%s", healthService)
	for k, v := range healthMap {
		healthService := fmt.Sprintf("%s %s\n", k, v)
		c.String(http.StatusOK, "%s", healthService)
	}
	//filehtml("health.html",healthMap)
	//c.JSON(http.StatusOK,healthMap)
	//c.File("./health.html")

}

func main() {
	// 创建一个默认的路由引擎
	// 默认使用了2个中间件Logger(),Recovery()
	r := gin.Default()
	// GET：请求方式；/hello：请求的路径
	// 一、当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("/", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})
	// 二、RESful
	r.GET("/healthchek", healthchek)
	// 四 URL参数
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run(":8000")
}
