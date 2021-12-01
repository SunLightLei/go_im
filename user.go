package main

import (
	"net"
	"strings"
)

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

	// 广播当前用户上线的消息
	this.server.BroadCast(this, "已上线！！")
}

// 用户下线
func (this *User) Offline() {
	// 用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户上线的消息
	this.server.BroadCast(this, "已下线。。")
}

// 给当前user对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线。。。\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式：rename|张三, 以|分割,rename下标0 张三为1
		newName := strings.Split(msg, "|")[1]
		// 判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已存在\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已经更新用户名: " + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式： to|userName|message
		// 1获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确， 请使用\"to|userName|message\".\n")
			return
		}
		// 2根据用户名，得到对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		// 3获取消息内容并发送出去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("无消息内容，请重发！\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对您说：" + content + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听用户的channel，有消息会发送给对方客户端
func (this *User) ListenMessage() {
	for {
		//channel的数据写进msg里
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
