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
"strings"
"time"
)


type SimpleRequest struct {
	Text SimpleRequestContent `json:"text"`
	Msgtype string `json:"msgtype"`
}

type SimpleRequestContent struct {
	Content string `json:"content"`
}
var UrlAddress = "https://oapi.dingtalk.com/robot/send?access_token=b291d1fc6728b6be13b786ec44809bb36425a33cec022c7782c98d19356f6f87"
var Secret = "SECf256c8ae405df3f958b613aa41e723b3a3aeaddff182e69547288dda2fa31a66"
func main() {
	// 设置请求头
	requestBody := SimpleRequest{
		Text: SimpleRequestContent{
			Content: "membership test",
		},
		Msgtype: "text",
	}
	reqBodyBox, _ := json.Marshal(requestBody)
	body := string(reqBodyBox)

	//  构建 签名
	//  把timestamp+"\n"+密钥当做签名字符串，使用HmacSHA256算法计算签名，然后进行Base64 encode，最后再把签名参数再进行urlEncode，得到最终的签名（需要使用UTF-8字符集）。
	timeStampNow := time.Now().UnixNano() / 1000000
	signStr :=fmt.Sprintf("%d\n%s", timeStampNow, Secret)

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
	fmt.Println("response: ", response)

	// 关闭 读取 reader
	defer response.Body.Close()

	// 读取内容
	all, _ := ioutil.ReadAll(response.Body)
	fmt.Println("all: ", string(all))
}
