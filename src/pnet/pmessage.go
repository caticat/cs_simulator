package pnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
)

type PMessage struct {
	// 网络数据
	command uint32        // 协议号
	proData proto.Message // 写数据用
	sliData []byte        // 读数据用

	// 计算用数据
	length uint32
}

func newPMessage() *PMessage {
	return &PMessage{}
}

func (this *PMessage) Command() uint32 { return this.command }

func (this *PMessage) Unmarshal(msg proto.Message) error {
	return proto.Unmarshal(this.sliData, msg)
}

func (this *PMessage) setCommand(command uint32) {
	this.command = command
}

func (this *PMessage) setMsg(msg proto.Message) {
	this.proData = msg
}

func (this *PMessage) setHead(data []byte) (length uint32, err error) {
	length = 0
	buffer := bytes.NewBuffer(data)
	//fmt.Println("收到头:", data)

	err = binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		fmt.Fprintf(os.Stderr, "出错了,读头长度:%v\n", err)
		return
	}

	err = binary.Read(buffer, binary.BigEndian, &this.command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "出错了,读command:%v\n", err)
		return
	}

	return
}

func (this *PMessage) setData(data []byte) {
	this.sliData = data
}

func (this *PMessage) serialize() []byte {
	buffer := bytes.NewBuffer([]byte{})

	sliData, err := proto.Marshal(this.proData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "出错了,proto.Marshal:%v\n", err)
		return nil
	}

	dataLength := uint32(len(sliData))
	err = binary.Write(buffer, binary.BigEndian, dataLength)
	if err != nil {
		fmt.Fprintf(os.Stderr, "出错了,写头长度:%v\n", err)
		return nil
	}

	err = binary.Write(buffer, binary.BigEndian, this.command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "出错了,写Command:%v\n", err)
		return nil
	}

	//fmt.Printf("发消息:协议号:%v,长度:%v\n", this.command, dataLength)

	buffer.Write(sliData)

	return buffer.Bytes()
}
