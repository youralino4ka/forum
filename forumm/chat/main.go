package main

import (
    "github.com/gorilla/websocket"
    "net/http"
)

var upgrader = websocket.Upgrader{}

func handleChat(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()
    // Логика обработки сообщений
}

func main() {
    http.HandleFunc("/chat", handleChat)
    http.ListenAndServe(":8082", nil)
}
