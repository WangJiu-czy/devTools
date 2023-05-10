package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

func initWs() *Manager {
	go WsManagerHandle()
	webSocketManager = &Manager{
		Group:       make(map[string]*Client, 100),
		Register:    make(chan *Client),
		UnRegister:  make(chan *Client),
		clientCount: 0,
	}
	return webSocketManager
}

var webSocketManager *Manager

type Client struct {
	Id          string
	Socket      *websocket.Conn
	Message     chan []byte
	MessageType int
	Exit        chan struct{}
}
type Manager struct {
	Group                map[string]*Client
	Lock                 sync.Mutex
	Register, UnRegister chan *Client
	clientCount          uint
}

const errCount = 3

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	//允许跨域请求
	return true
}}

/*
监听客户端注册
*/
func WsManagerHandle() {
	for {
		select {
		case client := <-webSocketManager.Register:
			webSocketManager.Lock.Lock()
			webSocketManager.Group[client.Id] = client
			webSocketManager.clientCount += 1
			log.Println(fmt.Sprintf("客户端注册: 客户端id为%s", client.Id))
			webSocketManager.Lock.Unlock()
		case client := <-webSocketManager.UnRegister:
			webSocketManager.Lock.Lock()
			if _, ok := webSocketManager.Group[client.Id]; ok {
				client.Exit <- struct{}{}
				close(client.Message)
				delete(webSocketManager.Group, client.Id)
				webSocketManager.clientCount -= 1
				log.Println(fmt.Sprintf("客户端注销: 客户端id为%s", client.Id))
			}
			webSocketManager.Lock.Unlock()
		}
	}
}
func (receiver *Client) Read() {
	count := 0
	for {
		mt, msg, err := receiver.Socket.ReadMessage()
		if err != nil {
			log.Println("连接错误")
			count++
			if count >= errCount {
				webSocketManager.UnRegister <- receiver
			}

			return
		}
		fmt.Println("读取执行")
		fmt.Println(msg)
		receiver.MessageType = mt
		//receiver.Message <- msg
		//消息类型switch{}
		//receiver.Message <- []byte("测试测试")
		select {
		case <-receiver.Exit:
			return
		default:

		}
	}
}
func (receiver *Client) Write() {
	count := 0
	for {

		err := receiver.Socket.WriteMessage(receiver.MessageType, <-receiver.Message)
		fmt.Println("写入执行")
		if err != nil {
			log.Println(receiver.Id + "----->写入错误")
			if count >= errCount {
				webSocketManager.UnRegister <- receiver
			}
			return
		}

		select {
		case <-receiver.Exit:
			return
		default:

		}
	}
}
