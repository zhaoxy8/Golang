package main

import (
	"sync"

	"Golong/aispace.com/tail_cpl/shell"
	"Golong/aispace.com/tail_cpl/taillog"
)

func main() {
	var wg sync.WaitGroup
	filename := "/opt/viv/type-server/logs/webservice-current.log"
	taillog.NewTailTask(filename)
	wg.Add(1)
	go shell.Init()
	wg.Wait()

	// a := say()
	// a()
	// a()
	// a = nil
	// // go a()
	// // go a()
	// // runtime.GC()
	// // fmt.Println(a())
	// // fmt.Println(say())
	// var b func()
	// b = nil
	// fmt.Println(b)
	// var c rune
	// var d []byte
	// c = 91
	// d = nil
	// fmt.Printf("%v\n", string(c))
	// fmt.Printf("%v\n", d)
}

// func say() func() {
// 	v := 0
// 	// b := "abc"
// 	sy := func() {
// 		v = v + 1
// 		fmt.Println(v)
// 	}
// 	// sy()
// 	return sy
// }
