package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	serverHost := "ws://127.0.0.1:8000/GET_websocket"
	var wsConn *websocket.Conn //websocket連線控制器
	go func() {
		for {
			if wsConn == nil {
				var err error
				wsConn, _, err = websocket.DefaultDialer.Dial(serverHost, nil)
				if err != nil {
					log.Printf("dial '%v' fail > '%v'\n", serverHost, err)
					time.Sleep(1 * time.Second)
					continue
				} else {
					log.Printf("suecess to connect websocket > '%v'\n", serverHost)
				}
			}

			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				log.Println(serverHost, " connection error: ", err)
				wsConn = nil //連現異常清空websocket連線控制器，在進行重連
				continue
			}
			log.Printf("%v receive data string: %s", serverHost, msg)
		}
	}()

	var sendMsg string
	for {
		fmt.Scanln(&sendMsg)
		sendJson := `{"Message":"` + sendMsg + `"}`
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(sendJson))
		if err != nil {
			log.Println("sent fail !")
		} else {
			log.Printf("sent %v\n", sendJson)
		}
		time.Sleep(1 * time.Second)
	}
	//wsConn.Close()
}
