package pnet

import (
	"fmt"
	"net"
)

type PServer struct {
	listenIP   string
	listenPort uint32

	listener    net.Listener
	chaClose    chan bool // TODO:PAN 这个东西实际没有用,待修正
	chaWaitRead chan *PMessageServer
	mapClient   map[string]*PServerClient // 所有的客户端连接
}

func NewPServer(ip string, port uint32) *PServer {
	return &PServer{
		listenIP:    ip,
		listenPort:  port,
		chaClose:    make(chan bool, 1),
		chaWaitRead: make(chan *PMessageServer, kPServerWaitReadLen),
		mapClient:   make(map[string]*PServerClient),
	}
}

func (this *PServer) Start() error {
	if err := this.listen(); err != nil {
		return err
	}

	go this.accept()

	return nil
}

func (this *PServer) Close() {
	this.chaClose <- true

	for _, serverClient := range this.mapClient {
		serverClient.Close()
	}

	this.listener.Close()
}

func (this *PServer) Read() chan *PMessageServer { return this.chaWaitRead }

func (this *PServer) onPServerClientClose(key string) {
	fmt.Printf("客户端[%v]断开连接", key)
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
