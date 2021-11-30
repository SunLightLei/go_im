package main

import (
	"fmt"
	"net"
)

// 声明Server
type Server struct {
	Ip   string
	Port int
}

// 创建一个server,传入ip和port,返回*Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("连接成功。。")
}

// 启动服务器
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen err: ", err.Error())
		return
	}
	// close socket listen
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
		}
		// do handler
		go this.Handler(conn)
	}
}
