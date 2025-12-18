package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"price-is-right-server/config"
	"price-is-right-server/game"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

var (
	rooms = make(map[string]*game.Room)
	mutex = &sync.Mutex{}
)

func getOrCreateRoom(roomID string) *game.Room {
	mutex.Lock()
	defer mutex.Unlock()

	if room, ok := rooms[roomID]; ok {
		return room
	}

	room := game.NewRoom(roomID)
	rooms[roomID] = room
	go room.Run()
	return room
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	roomID := r.URL.Query().Get("room")

	if name == "" || roomID == "" {
		http.Error(w, "Missing name or room", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	room := getOrCreateRoom(roomID)
	player := &game.Player{
		Name: name,
		Conn: conn,
	}

	room.Register <- player

	defer func() {
		room.Unregister <- player
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg game.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch msg.Type {
		case "START_GAME":
			room.StartGame <- player
		case "RESET_GAME":
			room.ResetGame <- player
		case "GUESS":
			// Payload should be float64
			if guessVal, ok := msg.Payload.(float64); ok {
				room.ProcessGuess <- game.GuessMessage{
					PlayerName: name,
					Guess:      guessVal,
				}
			} else {
                // JSON numbers are often float64 in Go interface{}, but let's be careful
                // If it comes as a string or something else.
                // Let's try to re-marshal and unmarshal the payload specifically if needed, 
                // or just cast.
                // A safer way is to define a struct for the payload.
                // Let's just assume the client sends a number.
                // Actually, `msg.Payload` is `interface{}`.
                // If I send `{"type": "GUESS", "payload": 123}`, it will be float64.
                room.ProcessGuess <- game.GuessMessage{
					PlayerName: name,
					Guess:      guessVal,
				}
			}
		}
	}
}

func main() {
	// Load configuration
	if err := config.LoadConfig("config.json"); err != nil {
		log.Fatal("Error loading config:", err)
	}

	http.HandleFunc("/ws", handleWebSocket)
	
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
