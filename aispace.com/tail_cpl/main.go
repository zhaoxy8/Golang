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
}
