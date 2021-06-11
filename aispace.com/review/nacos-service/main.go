/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)
//写入health.html文件中
func filehtml(name string,healthmap map[string]string){
	file, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	file.WriteString("# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.\n# TYPE go_gc_duration_seconds summary\n")
	for k, v := range  healthmap{
		fmt.Println(k, v)
		healthService := fmt.Sprintf("%s %s\n",k,v)
		file.WriteString(healthService)
	}
}
func healthchek(c *gin.Context) {
	healthMap := make(map[string]string, 100)
	//Membership
	//var m2service = [...]string{"membership-consumer-bizapi",
	//	                "membership-provider-marketactivity",
	//	                "membership-provider-memberservice",
	//	                "membership-common-auth",
	//	                "membership-provider-integrationcdp",
	//	                "membership-provider-points",
	//	                "membership-provider-pointsmall",
	//	                "membership-provider-integrationdmo",
	//	                "membership-consumer-pointsmall",
	//	                "membership-task",
	//	                "membership-provider-voucher",
	//	                "membership-provider-msg",
	//	                "membership-consumer-cpapi",
	//	                "membership-provider-cust",
	//	                "membership-consumer-integateapi",
	//                   }
    //Ecommerce
	var m2service = [...]string{"ecommerce-consumer-usercenter",
		"ecommerce-consumer-website-external",
		"ecommerce-provider-authority",
		"ecommerce-provider-common",
		"ecommerce-provider-payment",
		"ecommerce-provider-shoppingcart",
		"ecommerce-provider-customer",
		"ecommerce-provider-order",
		"ecommerce-consumer-website-mini",
		"ecommerce-provider-product",
		"ecommerce-provider-merchant",
		"ecommerce-provider-coupon",
		"ecommerce-consumer-web",
		"ecommerce-provider-vehicle",
		"ecommerce-consumer-website-nc",
		"ecommerce-provider-security",
		"ecommerce-consumer-campaign",
		"ecommerce-consumer-portal",
	}
	var groupname = "DEFAULT_GROUP"
	// 方法1：for循环遍历
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			//"m-nacos-int.bmw-emall.cn",
			"nacos-internal-dev.bmw-emall.cn",
			//"10.85.43.19",
			443,
			constant.WithScheme("https"),
			constant.WithContextPath("/nacos")),
	}
	//cc := constant.ClientConfig{
	//	NamespaceId:         "", //namespace id
	//	TimeoutMs:           5000,
	//	NotLoadCacheAtStart: true,
	//	RotateTime:          "1h",
	//	MaxAge:              3,
	//	LogLevel:            "info",
	//}
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(false),
		constant.WithLogDir(""),
		constant.WithCacheDir(""),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)
	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}
	for i := 0; i < len(m2service); i++ {
		//fmt.Println(m2service[i])
		instances := ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
			ServiceName: m2service[i],
			GroupName:   groupname,
		})
        for i :=0; i < len(instances); i++ {
        	//每个实例健康检查结果 0 健康 1 不健康
        	var healthResult string
        	var healthService string
        	var bytes []byte
        	//Membership
			//url := "/publicApi/health"
			//Ecommerce
			url := "/test/health"
			//fmt.Printf("%s,%s,%d \n", instances[i].ServiceName, instances[i].Ip, instances[i].Port)
			port := strconv.FormatUint(instances[i].Port, 10)
			url = "http://"+instances[i].Ip +":"+ port + url
			//fmt.Println(url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("http.Get err=", err)
			}
			//fmt.Println(resp.StatusCode,200)
			if resp.StatusCode != 200 {
				bytes = []byte("error")
			}else {
				var err error
				bytes, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("ioutil.ReadAll err=", err)
				}
			}

			fmt.Println(string(bytes))
			//container_health{container_name="membership-provider-cust",ip="10.85.46.102",port="10020"} 0
			//container_health{container_name="membership-provider-pointsmall",ip="",ip="10.85.46.103",port="20020"} 1
			healthService = fmt.Sprintf("container_health{container_name=\"%s\",ip=\"%s\",port=\"%d\"}",instances[i].ServiceName,instances[i].Ip, instances[i].Port)
			if string(bytes) == "SUCCESS"{
				healthResult = "1"
			}else {
				healthResult = "0"
			}
			healthMap[healthService] = healthResult

		}
	}

	for k, v := range  healthMap{
		healthService := fmt.Sprintf("%s %s\n",k,v)
		c.String(http.StatusOK, "%s",healthService)
	}
	//filehtml("health.html",healthMap)
	//c.JSON(http.StatusOK,healthMap)
	//c.File("./health.html")
}

func main(){
	// 创建一个默认的路由引擎
	// 默认使用了2个中间件Logger(),Recovery()
	r := gin.Default()
	// GET：请求方式；/hello：请求的路径
	// 一、当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("/hello", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})
	// 二、RESful
	r.GET("/healthchek",healthchek)
	// 四 URL参数
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run(":8000")
}