// 协议编译
// protoc3 test.proto --go_out=../pcmd

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"pnet"
)

var (
	gServer = pnet.NewPServer("0.0.0.0", 10001)
)

func main() {
	fmt.Println("服务器启动")

	//test()
	//return

	chaClose := make(chan bool, 1)
	if err := gServer.Start(chaClose); err != nil {
		fmt.Println("启动服务器失败:", err)
		return
	}
	defer gServer.Close()

	loop(chaClose)

	fmt.Println("服务器停止")
}

func test() {
	d := uint32(12)
	a := bytes.NewBuffer([]byte{})
	binary.Write(a, binary.BigEndian, d)
	fmt.Println("数据:", a.Bytes())
	c := uint32(0)
	b := bytes.NewBuffer(a.Bytes())
	binary.Read(b, binary.BigEndian, &c)
	fmt.Println("读取:", c)
}
