package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
)

func TestClient_Read(t *testing.T) {
	_ = InitDefaultWs()
	app := gin.Default()
	app.GET("/ws", WsCl)
	app.Run(":8080")
}

func WsCl(c *gin.Context) {
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}

	var client *tClient = &tClient{Client: &Client{}}

	client.Id = c.Query("userId")
	client.Socket = ws
	client.Message = make(chan []byte, 1024)

	webSocketManager.Register <- client.Client
	go client.Read()
	go client.Write()

}

type tClient struct {
	*Client
}

func (receiver *tClient) Read() {
	count := 0
	for {
		mt, msg, err := receiver.Socket.ReadMessage()
		if err != nil {
			log.Println("连接错误")
			count++
			if count >= errCount {
				webSocketManager.UnRegister <- receiver.Client
			}

			return
		}
		fmt.Println("tClient读取执行")
		fmt.Println(msg)
		receiver.MessageType = mt
		select {
		case <-receiver.Exit:
			return
		default:

		}
	}
}
func (receiver *tClient) Write() {
	count := 0
	for {

		err := receiver.Socket.WriteMessage(receiver.MessageType, <-receiver.Message)
		fmt.Println("tClient写入执行")
		if err != nil {
			log.Println(receiver.Id + "----->写入错误")
			if count >= errCount {
				webSocketManager.UnRegister <- receiver.Client
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
