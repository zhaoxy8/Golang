package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)
var personSlice []Person
var requestMap map[string]uint8
var remoteMap map[string]uint8

type Person struct {
	Request   string  `json:"request"`
	Status    string  `json:"status"`
	Remote_addr string  `json:"remote_addr"`
}



func main() {
	var code int = 200
	requestMap = make(map[string]uint8)
	remoteMap = make(map[string]uint8)
	file, err := os.Open("./access.log")
	if err != nil{
		fmt.Printf("open file failed,err: ",err)
		return
	}
	defer  file.Close()
	reader := bufio.NewReader(file)
	personSlice = make([]Person,0,100)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF{
			if len(line) != 0 {
				var p2 Person
				json.Unmarshal([]byte(line),&p2)
				personSlice = append(personSlice,p2)
			}
			fmt.Println("文件读取完毕")
			break
		}
		if err != nil{
			fmt.Println("read file failed, err:",err)
			return
		}
		var p2 Person
		json.Unmarshal([]byte(line),&p2)
		//fmt.Println(p2)
		personSlice = append(personSlice,p2)
		//fmt.Println(personSlice)
	}
	for index,p := range personSlice{
		stat,_ := strconv.Atoi(p.Status)
		if stat == code {
			_ ,ok := remoteMap[p.Request]
			if ok {
				requestMap[p.Request] += 1
			}else{
				requestMap[p.Request] = 1
			}
			fmt.Printf("第%d条: %v\n",index,p.Remote_addr)
		}
	}
	for k,v := range remoteMap{
		fmt.Println("====",k,v)
	}
}
