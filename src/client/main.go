package main

import (
	"fmt"
	"pnet"
)

var (
	gClient = pnet.NewPClient("127.0.0.1", 10001)
)

func main() {
	fmt.Println("开始")

	if err := gClient.Start(); err != nil {
		fmt.Println("dail出错了:%v", err)
		return
	}
	defer gClient.Close()

	loop()

	fmt.Println("结束")
}
