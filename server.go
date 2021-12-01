package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// 定义Server
type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	message chan string
}

// 创建一个server,传入ip和port,返回*Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		message:   make(chan string),
	}
	return server
}

// 监听message广播消息 channel的goroutine，有消息就发送给所有在线的用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.message

		//将msg发送给所有在线的用户
		this.mapLock.Lock()
		// cli表示用户，
		for _, cli := range this.OnlineMap {
			// 将msg 写进用户的channel里
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息方法
func (this *Server) BroadCast(user *User, msg string) {
	// 定义广播消息内容和格式
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	// 将广播消息传给 message(这是个channel)
	this.message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("连接成功。。")
	user := NewUser(conn, this)

	//用户上线，将该用户加入OnlineMap中
	//this.mapLock.Lock()
	//this.OnlineMap[user.Name] = user
	//this.mapLock.Unlock()
	user.Online()

	// 广播当前用户上线的消息
	//this.BroadCast(user, "已上线。。")

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		// 建一个buffer，大小为4096
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			// 如果没有消息过来，表示对端的客户端已经下线
			if n == 0 {
				//this.BroadCast(user, "已下线！！")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}

			// 提取消息(去除'\n')
			msg := string(buf[:n-1])

			// 将得到的msg进行广播
			//this.BroadCast(user, msg)

			// 用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息，代表当前用户是活跃的
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，重置定时器
		case <-time.After(time.Second * 10):
			// 超时，关闭当前用户
			user.SendMsg("您已下线。。。")

			// 销毁用户的资源
			close(user.C)

			// 关闭连接
			conn.Close()

			// 退出当前的handler
			return
		}
	}
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

	// 启动监听message的 goroutine
	go this.ListenMessager()

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
