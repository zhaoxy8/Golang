package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter,r *http.Request){
	tmpl , err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		fmt.Println(err)
	}
	name := "泰迪"
	err = tmpl.ExecuteTemplate(w,"index.tmpl",name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/index",index)
	err := http.ListenAndServe(":9090",nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
