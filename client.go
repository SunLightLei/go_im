package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //记录用户选择的聊天模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	client.conn = conn

	// 返回对象
	return client
}

// 菜单
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("====请输入合法范围内的数据====")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println("您选择的是公聊模式。。。")
			break
		case 2:
			// 私聊模式
			fmt.Println("您选择的是私聊模式。。。")
			break
		case 3:
			// 重命名
			fmt.Println("您选择的是重命名模式。。。")
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {

	// 命令行解析
	flag.Parse()

	//client := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败......")
		return
	}

	fmt.Println("连接服务器成功......")

	// 启动客户端业务
	select {}
}
