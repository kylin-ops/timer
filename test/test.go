package main

import (
	"fmt"
	"github.com/kylin-ops/timer"
	"time"
)

func main() {
	t := timer.NewTimer()
	fmt.Println("start aaaa")
	t.Add("test", time.Second*5, func() {
		panic("aaa")
		fmt.Println("aaaa")
	})
	fmt.Println("start bbbb")
	t.Add("test1", time.Second*2, func() {
		fmt.Println("bbbb")
	})
	time.Sleep(time.Second * 6)
	t.Del("test1")
	time.Sleep(time.Second * 60)
	fmt.Println("end")
}
