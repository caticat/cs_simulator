package pnet

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"os"
	"sync"
)

type PServerClient struct {
	conn         net.Conn
	chaWaitRead  chan *PMessageServer
	chaWaitWrite chan *PMessage

	// 关闭相关处理(这里写的好复杂-_-||)
	chaClose   chan bool
	funOnClose func(string)
	waiClient  *sync.WaitGroup
	mutClose   sync.Mutex
	isClosing  bool
}

func newPServerClient(server *PServer, conn net.Conn) *PServerClient {
	return &PServerClient{
		conn:         conn,
		chaWaitRead:  server.Read(),
		chaWaitWrite: make(chan *PMessage, kPServerClientWaitSendLen),
		chaClose:     make(chan bool, 1),
		funOnClose:   server.onPServerClientClose,
		waiClient:    server.waiClient,
		isClosing:    false,
	}
}

func (this *PServerClient) Send(cmd uint32, msg proto.Message) {
	pmsg := newPMessage()
	pmsg.setCommand(cmd)
	pmsg.setMsg(msg)
	this.chaWaitWrite <- pmsg
}

func (this *PServerClient) Close() {
	this.mutClose.Lock()
	defer this.mutClose.Unlock()

	if this.isClosing {
		return
	}
	this.isClosing = true
	this.chaClose <- true
}

func (this *PServerClient) start() {
	this.waiClient.Add(2)
	go this.loopRead()
	go this.loopSend()
}

func (this *PServerClient) getClientAddr() string {
	return this.conn.RemoteAddr().String()
}

func (this *PServerClient) loopRead() {
	defer this.waiClient.Done()
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
	defer this.waiClient.Done()
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
