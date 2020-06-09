package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Name string
	Age int
	Gender string
}

func sayhello(w http.ResponseWriter,r *http.Request){
	tmpl,err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	u1 := User{
		Name: "泰迪",
		Age: 29,
		Gender: "男",
	}
	m1 := map[string]interface{}{
		"host":"10.0.0.1",
		"result":"SUCCESS",
	}
	hobbySlice := []string{
		"唱歌",
		"水机",
		"工作",
	}
	tmpl.Execute(w,map[string]interface{}{
		"u1":u1,
		"m1":m1,
		"hobby":hobbySlice,
	})
}

func main() {
	http.HandleFunc("/",sayhello)
	err := http.ListenAndServe(":9090",nil)
	if err != nil{
		fmt.Printf("HTTP server failed,err:%v\n", err)
		return
	}
}
