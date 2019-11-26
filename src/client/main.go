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

	chaClose := make(chan bool, 1)
	if err := gClient.Start(chaClose); err != nil {
		fmt.Println("dail出错了:%v", err)
		return
	}
	defer gClient.Close()

	loop(chaClose)

	fmt.Println("结束")
}
