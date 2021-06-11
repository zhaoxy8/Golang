package main

import (
"fmt"

"github.com/nacos-group/nacos-sdk-go/clients"
"github.com/nacos-group/nacos-sdk-go/common/constant"
"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	//sc := []constant.ServerConfig{
	//	{
	//		IpAddr: "m-nacos-dev.bmw-emall.cn",
	//		Port:   443,
	//	},
	//}
	//or a more graceful way to create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			"m-nacos-dev.bmw-emall.cn",
			443,
			constant.WithScheme("https"),
			constant.WithContextPath("/nacos")),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	//or a more graceful way to create ClientConfig
	//_ = *constant.NewClientConfig(
	//	constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"),
	//	constant.WithTimeoutMs(5000),
	//	constant.WithNotLoadCacheAtStart(true),
	//	constant.WithLogDir("/tmp/nacos/log"),
	//	constant.WithCacheDir("/tmp/nacos/cache"),
	//	constant.WithRotateTime("1h"),
	//	constant.WithMaxAge(3),
	//	constant.WithLogLevel("debug"),
	//)

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//publish config
	//config key=dataId+group+namespaceId
	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId:  "test-data",
	//	Group:   "test-group",
	//	Content: "hello world!",
	//})
	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId:  "test-data-2",
	//	Group:   "test-group",
	//	Content: "hello world!",
	//})
	//if err != nil {
	//	fmt.Printf("PublishConfig err:%+v \n", err)
	//}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "membership-gateway-dev.yaml",
		Group:  "V1_GROUP",
	})
	fmt.Println("GetConfig,config :" + content)

	//Listen config change,key=dataId+group+namespaceId.
	//err = client.ListenConfig(vo.ConfigParam{
	//	DataId: "test-data",
	//	Group:  "test-group",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
	//	},
	//})

	//err = client.ListenConfig(vo.ConfigParam{
	//	DataId: "test-data-2",
	//	Group:  "test-group",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
	//	},
	//})

	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId:  "test-data",
	//	Group:   "test-group",
	//	Content: "test-listen",
	//})

	//time.Sleep(2 * time.Second)

	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId:  "test-data-2",
	//	Group:   "test-group",
	//	Content: "test-listen",
	//})

	//time.Sleep(2 * time.Second)

	////cancel config change
	//err = client.CancelListenConfig(vo.ConfigParam{
	//	DataId: "test-data",
	//	Group:  "test-group",
	//})

	//time.Sleep(2 * time.Second)
	//_, err = client.DeleteConfig(vo.ConfigParam{
	//	DataId: "test-data",
	//	Group:  "test-group",
	//})
	//time.Sleep(5 * time.Second)

	//searchPage, _ := client.SearchConfig(vo.SearchConfigParm{
	//	Search:   "blur",
	//	DataId:   "",
	//	Group:    "",
	//	PageNo:   1,
	//	PageSize: 10,
	//})
	//fmt.Printf("Search config:%+v \n", searchPage)
}
