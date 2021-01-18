package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clientList = make(map[*websocket.Conn]struct{}) // 此map負責儲存WebSocket客戶端所發出 GET請求文本
var wsUpgrader = websocket.Upgrader{}               // 負責WebSocket更新與儲存Request與Response
var broadcastCh = make(chan JSONMsg)                // 為了加快伺服器處理，建立JSON訊息管道存放
type JSONMsg struct {                               // gorilla/websocket 內建JSON解析方便瀏覽器傳遞
	ClentMsg string `json:"messagE"` //gorilla ReadJSON不分大小寫,WriteJSON
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./Public"))) //將 ./Public資料夾內檔案公開給客戶端瀏覽
	http.HandleFunc("/GET_websocket", handleWS)             //接收自客戶端GET請求(大小寫區分)，並將請求指定函式處理

	err := http.ListenAndServe(":8000", nil) //localhost:8080 設定Port
	if err != nil {
		log.Fatal("http.ListenAndServe() error>", err)
	}

	close(broadcastCh)
	return

}

func handleWS(input_httpRW http.ResponseWriter, input_httpRQ *http.Request) {

	wsConn, err := wsUpgrader.Upgrade(input_httpRW, input_httpRQ, nil)
	if err != nil {
		log.Fatal("wsUpgrader.Upgrade error>", err)
	}
	defer wsConn.Close()

	vs_ClientIPnPort := (*wsConn).RemoteAddr()
	log.Printf("Got connect request from : %v\n", vs_ClientIPnPort) //紀錄客戶端IP&Port

	clientList[wsConn] = struct{}{} //初始化連線數

	go func() {
		for {
			//利用管道阻塞特性，該函式僅會在<-broadcastCh有元素狀況下才會往下執行
			message := <-broadcastCh
			log.Printf("Got from client: %v", message) //紀錄所接收用戶傳遞訊息

			for conn := range clientList {
				err := conn.WriteJSON(message) //將訊息回傳給客戶端
				if err != nil {
					log.Printf("conn.WriteJSON error> %v", err)
					conn.Close()
					delete(clientList, conn)
				}
			}
		}
	}()

	for {
		var ClientMsg JSONMsg
		err := wsConn.ReadJSON(&ClientMsg) //接收自客戶端傳遞訊息
		if err != nil {
			log.Printf("wsConn.ReadJSON error [%v]> %v", vs_ClientIPnPort, err)
			delete(clientList, wsConn) //將異常客戶端由清單中刪除
			break
		}
		broadcastCh <- ClientMsg //將所接收客戶端訊息依序塞入管道
	}
}
