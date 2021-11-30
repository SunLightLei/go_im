package main

import "net"

// 定义User
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前用户的channel消息的goroutine
	go user.ListenMessage()
	return user
}

// 用户上线
func (this *User) Online() {
	// 用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线
	this.server.BroadCast(this, "已上线！！")
}

// 用户下线
func (this *User) Offline() {
	// 用户下线，将用户从onlineMap中去掉
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户下线消息
	this.server.BroadCast(this, "已下线。。")
}

// 用户处理消息
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// 监听用户的channel，有消息会发送给对方客户端
func (this *User) ListenMessage() {
	for {
		//channel的数据写进msg里
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
