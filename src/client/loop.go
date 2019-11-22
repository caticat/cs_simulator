package main

import (
	"fmt"
	"pcmd"
	"time"
)

func loop() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case msg := <-gClient.Read():
			pmsg := &pcmd.ProtoCommandTest{}
			err := msg.Unmarshal(pmsg)
			if err != nil {
				fmt.Printf("客户端解析协议出错:%v\n", err)
				return
			}
			fmt.Printf("收到消息:协议号:%+v,协议内容:id:%v,data:%v\n", msg.Command(), pmsg.GetId(), pmsg.GetData())
		case <-ticker.C:
			//fmt.Println("定时器一跳")

			gClient.Send(uint32(pcmd.EProtoCommand_EProtoCommandTest), &pcmd.ProtoCommandTest{
				Id:   198,
				Data: "啊啊啊123abc",
			})
			//fmt.Println("发送消息了!")
		}
	}
}
