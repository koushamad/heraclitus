package service

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectToWebSocket(host string, port int) {

	conn, _, err := websocket.DefaultDialer.Dial("ws://"+host+":"+strconv.Itoa(port)+"/ws", nil)
	if err != nil {
		log.Fatal("Dial:", err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to %s:%d\n", host, port)

	for {
		_, request, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		reqReader := bufio.NewReader(bytes.NewReader(request))
		req, err := http.ReadRequest(reqReader)
		if err != nil {
			log.Println("ReadRequest:", err)
			continue
		}

		log.Println(req.Method, req.URL)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Do:", err)
			continue
		}

		dump, err := httputil.DumpResponse(resp, true)
		resp.Body.Close()
		if err != nil {
			log.Println("DumpResponse:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, dump)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
