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
	"net/url"
	"strings"
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

// Define the mapping from domain names to IP addresses and ports.
var domainToIPPort = map[string]string{
	"test.com": "http://127.0.0.1:3000",
}

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
		log.Println("No WebSocket clients connected")
		mu.Unlock()
		return
	}

	// Select a WebSocket client at random.
	rand.Seed(time.Now().UnixNano())
	conn := conns[rand.Intn(len(conns))]
	mu.Unlock()

	// Redirect the request to the corresponding IP and port.
	host := strings.Split(r.Host, ":")[0] // Get the domain name.
	if ipPort, ok := domainToIPPort[host]; ok {
		u, err := url.Parse(ipPort + r.URL.Path)

		if err != nil {
			http.Error(w, "Error redirecting request", http.StatusInternalServerError)
			log.Println("Error parsing URL:", err)
			return
		}

		r.URL = u
	}

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, "Error dumping request", http.StatusInternalServerError)
		log.Println("Error dumping request:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, dump)
	if err != nil {
		http.Error(w, "Error writing to WebSocket client", http.StatusInternalServerError)
		log.Println("Error writing to WebSocket client:", err)
		return
	}

	_, response, err := conn.ReadMessage()
	if err != nil {
		http.Error(w, "Error reading response from WebSocket client", http.StatusInternalServerError)
		log.Println("Error reading response from WebSocket client:", err)
		return
	}

	respReader := bufio.NewReader(bytes.NewReader(response))
	resp, err := http.ReadResponse(respReader, r)
	if err != nil {
		http.Error(w, "Error reading HTTP response", http.StatusInternalServerError)
		log.Println("Error reading HTTP response:", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
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
