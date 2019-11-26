package main

import (
	"fmt"
	"pcmd"
	"time"
)

func loop(chaClose chan bool) {
	ticker := time.NewTicker(time.Second)
	//i := 0
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
		case <-chaClose:
			fmt.Println("程序结束")
			return
		case <-ticker.C:
			//i++
			//if i == 10 {
			//	gServer.Close()
			//}
			//fmt.Println("定时器一跳")
		}
	}
}
