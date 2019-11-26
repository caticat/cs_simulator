package pnet

import (
	"fmt"
	"net"
	"sync"
)

type PServer struct {
	listenIP   string
	listenPort uint32

	listener    net.Listener
	chaWaitRead chan *PMessageServer
	mapClient   map[string]*PServerClient // 所有的客户端连接
	waiClient   *sync.WaitGroup
	chaClose    chan bool
}

func NewPServer(ip string, port uint32) *PServer {
	return &PServer{
		listenIP:    ip,
		listenPort:  port,
		chaWaitRead: make(chan *PMessageServer, kPServerWaitReadLen),
		mapClient:   make(map[string]*PServerClient),
		waiClient:   &sync.WaitGroup{},
	}
}

func (this *PServer) Start(chaClose chan bool) error {
	this.chaClose = chaClose
	if err := this.listen(); err != nil {
		return err
	}

	go this.accept()

	return nil
}

func (this *PServer) Close() {
	fmt.Println("开始关闭连接")
	for _, serverClient := range this.mapClient {
		serverClient.Close()
	}

	this.waiClient.Wait()

	this.listener.Close()

	this.chaClose <- true

	fmt.Println("连接关闭完成")
}

func (this *PServer) Read() chan *PMessageServer { return this.chaWaitRead }

func (this *PServer) onPServerClientClose(key string) {
	fmt.Printf("客户端[%v]断开连接\n", key)
	delete(this.mapClient, key)
}

func (this *PServer) listen() error {
	listener, err := net.Listen("tcp", this.getHost())
	if err != nil {
		fmt.Println("启动监听失败:", err)
		return err
	}

	this.listener = listener

	return nil
}

func (this *PServer) accept() {
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			fmt.Println("accept出错:", err)
			return
		}
		clientHost := conn.RemoteAddr().String()
		serverClient := newPServerClient(this, conn)
		this.mapClient[clientHost] = serverClient
		serverClient.start()
		fmt.Printf("客户端[%v]连接\n", clientHost)
	}
}

func (this *PServer) getHost() string {
	return fmt.Sprintf("%v:%v", this.listenIP, this.listenPort)
}
