package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type UserInfo struct {
	Name string
	Gender string
	Age int
}

func tmplDemo(w http.ResponseWriter,r *http.Request){
	tmpl,err := template.ParseFiles("./hello.tmpl","./ul.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	user := UserInfo{
		Name: "",
		Gender: "",
		Age: 18,
	}
	tmpl.Execute(w,user)
}

func main() {
http.HandleFunc("/tmpl",tmplDemo)
err := http.ListenAndServe(":9090",nil)
if err != nil{
	fmt.Printf("HTTP server failed,err:%v\n", err)
	return
}
}
