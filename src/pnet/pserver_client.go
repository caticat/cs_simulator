package pnet

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"os"
)

type PServerClient struct {
	conn         net.Conn
	chaWaitRead  chan *PMessageServer
	chaWaitWrite chan *PMessage
	chaClose     chan bool
	funOnClose   func(string)
}

func newPServerClient(server *PServer, conn net.Conn) *PServerClient {
	return &PServerClient{
		conn:         conn,
		chaWaitRead:  server.Read(),
		chaWaitWrite: make(chan *PMessage, kPServerClientWaitSendLen),
		chaClose:     make(chan bool, 1),
		funOnClose:   server.onPServerClientClose,
	}
}

func (this *PServerClient) Send(cmd uint32, msg proto.Message) {
	pmsg := newPMessage()
	pmsg.setCommand(cmd)
	pmsg.setMsg(msg)
	this.chaWaitWrite <- pmsg
}

func (this *PServerClient) Close() {
	this.chaClose <- true
}

func (this *PServerClient) start() {
	go this.loopRead()
	go this.loopSend()
}

func (this *PServerClient) getClientAddr() string {
	return this.conn.RemoteAddr().String()
}

func (this *PServerClient) loopRead() {
	head := make([]byte, kPMessageHeadLen, kPMessageHeadLen)
	for {
		_, err := io.ReadFull(this.conn, head)
		if err != nil {
			this.Close()
			// MARK: 这里的错误无视,只当成是连接断开处理
			//fmt.Fprintf(os.Stderr, "出错了,读取头数据出错:%v\n", err)
			return
		}

		msg := newPMessageServer(this)
		length, err := msg.setHead(head)
		if err != nil {
			this.Close()
			fmt.Fprintf(os.Stderr, "出错了,解析头数据出错:%v\n", err)
			return
		}

		//fmt.Printf("消息协议号:[%v],长度:[%v]\n", msg.command, length)
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

func (this *PServerClient) loopSend() {
	for {
		select {
		case msg := <-this.chaWaitWrite:
			this.send(msg)
		case <-this.chaClose:
			this.close()
			return
		}
	}
}

func (this *PServerClient) send(msg *PMessage) {
	this.conn.Write(msg.serialize())
}

func (this *PServerClient) close() {
	this.funOnClose(this.getClientAddr())

	close(this.chaClose)
	this.conn.Close()
}
