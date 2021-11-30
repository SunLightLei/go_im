package main

import "net"

// 定义User
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	// 启动监听当前用户的channel消息的goroutine
	go user.ListenMessage()
	return user
}

// 监听用户的channel，有消息会发送给对方客户端
func (this *User) ListenMessage() {
	for {
		//channel的数据写进msg里
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
