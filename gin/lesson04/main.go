package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func sayhello(w http.ResponseWriter,r *http.Request){
	tmpl,err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	tmpl.Execute(w,"王磊")
}

func main() {
http.HandleFunc("/",sayhello)
err := http.ListenAndServe(":9090",nil)
if err != nil{
	fmt.Printf("HTTP server failed,err:%v\n", err)
	return
}
}
