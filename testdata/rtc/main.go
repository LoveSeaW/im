package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var offerClient *websocket.Conn
var answerClient *websocket.Conn

func checkStart() {
	if offerClient != nil && answerClient != nil {
		offerClient.WriteJSON(map[string]string{
			"type": "create_offer",
		})
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 鉴权 true表示放行，false表示拦截
			return true
		},
	}
	fmt.Println("建立ws连接")
	// 将http升级至websocket
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	for {
		var obj map[string]any
		err := conn.ReadJSON(&obj)
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		switch obj["type"] {
		case "connect":
			if offerClient == nil {
				offerClient = conn
				conn.WriteJSON(map[string]interface{}{
					"type":    "connect",
					"code":    200,
					"message": "connect success",
				})
				checkStart() // 是否可以准备开始
			} else if answerClient == nil {
				answerClient = conn
				conn.WriteJSON(map[string]interface{}{
					"type":    "connect",
					"code":    200,
					"message": "connect success",
				})
				checkStart()
			} else {
				conn.WriteJSON(map[string]interface{}{
					"type":    "connect",
					"code":    -1,
					"message": "connect failed",
				})
				conn.Close()
			}
		case "offer": // offer
			if answerClient != nil {

				fmt.Println("offer:", obj["type"])
				answerClient.WriteJSON(obj)
			}
		case "answer": // 应答
			if offerClient != nil {

				fmt.Println("answer:", obj["type"])
				offerClient.WriteJSON(obj)
			}
		case "offer_ice":
			if answerClient != nil {

				fmt.Println("offer_ice:", obj["type"])
				answerClient.WriteJSON(obj)
			}
		case "answer_ice":
			if offerClient != nil {

				fmt.Println("answer_ice:", obj["type"])
				offerClient.WriteJSON(obj)
			}
		}
	}

	if conn == offerClient {
		log.Println("remove offerClient")
		offerClient = nil
	} else if conn == answerClient {
		log.Println("remove answerClient")
		answerClient = nil
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	byteData, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(byteData)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server running on :9004")
	log.Fatal(http.ListenAndServe(":9004", nil))
}
