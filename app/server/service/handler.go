package service

import (
	"bufio"
	"bytes"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	mu    sync.Mutex
	conns []*websocket.Conn
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	mu.Lock()
	conns = append(conns, conn)
	mu.Unlock()

	log.Printf("WebSocket client connected. Total clients: %d\n", len(conns))
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if len(conns) == 0 {
		http.Error(w, "No WebSocket clients connected", http.StatusInternalServerError)
		mu.Unlock()
		return
	}

	// Select a WebSocket client at random.
	rand.Seed(time.Now().UnixNano())
	conn := conns[rand.Intn(len(conns))]
	mu.Unlock()

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, "Error dumping request", http.StatusInternalServerError)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, dump)
	if err != nil {
		http.Error(w, "Error writing to WebSocket client", http.StatusInternalServerError)
		return
	}

	_, response, err := conn.ReadMessage()
	if err != nil {
		http.Error(w, "Error reading response from WebSocket client", http.StatusInternalServerError)
		return
	}

	respReader := bufio.NewReader(bytes.NewReader(response))
	resp, err := http.ReadResponse(respReader, r)
	if err != nil {
		http.Error(w, "Error reading HTTP response", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
