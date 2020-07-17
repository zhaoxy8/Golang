package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)
/*
  处理nignx日志，找到httpStatus=200的数据
  统计URI地址访问最多的top10
  统计upstream地址访问最多的top10
 */
var PersonSlice []Person
var RequestMap map[string]uint8
var RemoteMap map[string]uint8

type Person struct {
	Request     string `json:"request"`
	Status      string `json:"status"`
	Remote_addr string `json:"remote_addr"`
}

//对map的value进行排序
//要对golang map按照value进行排序，思路是直接不用map，用struct存放key和value，实现sort接口，就可以调用sort.Sort进行排序了。
type Pair struct {
	Key   string
	Value uint8
}

// A slice of Pairs that implements sort.Interface to sort by Value.
//自定义类型
type PairList []Pair
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]uint8) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	fmt.Println(p)
	return p
}

func main() {
	var code int = 200
	RequestMap = make(map[string]uint8)
	RemoteMap = make(map[string]uint8)
	file, err := os.Open("./access.log")
	if err != nil {
		fmt.Printf("open file failed,err: ", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	PersonSlice = make([]Person, 0, 100)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(line) != 0 {
				var p2 Person
				json.Unmarshal([]byte(line), &p2)
				PersonSlice = append(PersonSlice, p2)
			}
			fmt.Println("文件读取完毕")
			break
		}
		if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		var p2 Person
		json.Unmarshal([]byte(line), &p2)
		//fmt.Println(p2)
		PersonSlice = append(PersonSlice, p2)
		//fmt.Println(personSlice)
	}
	for index, p := range PersonSlice {
		stat, _ := strconv.Atoi(p.Status)
		if stat == code {
			_, ok := RequestMap[p.Request]
			if ok {
				RequestMap[p.Request] += 1
			} else {
				RequestMap[p.Request] = 1
			}
			fmt.Printf("第%d条: %v\n", index, p.Request)

			_, ok = RemoteMap[p.Remote_addr]
			if ok {
				RemoteMap[p.Remote_addr] += 1
			} else {
				RemoteMap[p.Remote_addr] = 1
			}
		}
	}
	p1 := sortMapByValue(RequestMap)
	p2 := sortMapByValue(RemoteMap)
	for i := 0; i < 3; i++ {
		fmt.Println(p1[i])
	}
	for i := 0; i < 3; i++ {
		fmt.Println(p2[i])
	}
	//for k, v := range RequestMap {
	//	fmt.Printf("%s====%d\n", k, v)
	//}
}
