package pnet

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"os"
)

type PClient struct {
	remoteIP   string
	remotePort uint32

	conn        net.Conn
	chaClose    chan bool
	chaWaitSend chan *PMessage
	chaWaitRead chan *PMessage
}

func NewPClient(ip string, port uint32) *PClient {
	return &PClient{
		remoteIP:    ip,
		remotePort:  port,
		chaClose:    make(chan bool, 1),
		chaWaitSend: make(chan *PMessage, kPClientWaitSendLen),
		chaWaitRead: make(chan *PMessage, kPClientWaitReadLen),
	}
}

func (this *PClient) Close() {
	this.chaClose <- true
}

func (this *PClient) Start() error {
	if err := this.dail(); err != nil {
		return err
	}

	go this.loopRead()
	go this.loopSend()

	return nil
}

func (this *PClient) Send(cmd uint32, msg proto.Message) {
	pmsg := newPMessage()
	pmsg.setCommand(cmd)
	pmsg.setMsg(msg)
	this.chaWaitSend <- pmsg
}

func (this *PClient) Read() chan *PMessage {
	return this.chaWaitRead
}

func (this *PClient) dail() error {
	conn, err := net.Dial("tcp", this.getHost())
	if err != nil {
		return err
	}
	this.conn = conn

	return nil
}

func (this *PClient) loopRead() {
	head := make([]byte, kPMessageHeadLen, kPMessageHeadLen)
	for {
		//fmt.Println("收到消息了!!!")
		_, err := io.ReadFull(this.conn, head)
		if err != nil {
			this.Close()
			// MARK: 这里的错误无视,只当成是连接断开处理
			//fmt.Fprintf(os.Stderr, "出错了,读取头数据出错:%v\n", err)
			return
		}

		msg := newPMessage()
		length, err := msg.setHead(head)
		if err != nil {
			this.Close()
			fmt.Fprintf(os.Stderr, "出错了,解析头数据出错:%v\n", err)
			return
		}

		data := make([]byte, length, length)
		_, err = io.ReadFull(this.conn, data)
		if err != nil {
			this.Close()
			fmt.Fprintf(os.Stderr, "出错了,读body数据出错:%v\n", err)
			return
		}
		msg.setData(data)

		this.chaWaitRead <- msg
	}
}

func (this *PClient) loopSend() {
	for {
		select {
		case msg := <-this.chaWaitSend:
			this.send(msg)
		case <-this.chaClose:
			this.close()
			return
		}
	}
}

func (this *PClient) close() {
	close(this.chaClose)
	this.conn.Close()
}

func (this *PClient) send(msg *PMessage) {
	this.conn.Write(msg.serialize())
}

func (this *PClient) getHost() string {
	return fmt.Sprintf("%v:%v", this.remoteIP, this.remotePort)
}
