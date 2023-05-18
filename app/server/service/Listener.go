package service

import (
	"log"
	"net/http"
	"strconv"
)

func Listen(port int) {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", handleHTTPRequest)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
