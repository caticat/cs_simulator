package main

import (
	"fmt"
	"pcmd"
	"time"
)

func loop() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case msg := <-gServer.Read():
			pmsg := &pcmd.ProtoCommandTest{}
			err := msg.Unmarshal(pmsg)
			if err != nil {
				fmt.Println("收消息出错了:", err)
				continue
			}
			fmt.Printf("收到消息:客户端:[%+v], 协议号:[%v],消息内容:[id:%v,data:%v]\n", msg.GetClientAddr(), msg.Command(), pmsg.GetId(), pmsg.GetData())
			msg.GetClient().Send(uint32(pcmd.EProtoCommand_EProtoCommandTest), &pcmd.ProtoCommandTest{
				Id:   998,
				Data: "啦啦啦,123,abc",
			})
		case <-ticker.C:
			//fmt.Println("定时器一跳")
		}
	}
}
