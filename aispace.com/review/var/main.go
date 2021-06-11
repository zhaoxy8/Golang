package main

import (
	"fmt"
	"time"
)

func main(){
	now := time.Now()
	//
	//2021-05-14T11:05:32
	fmt.Println(now.Format("2006-01-02T15:04:05.000"))
	var a int = 10
	fmt.Printf("%d \n",a)
	fmt.Printf("%b \n",a)
	s1 := "hello"
	fmt.Printf("%s \n",s1)
}