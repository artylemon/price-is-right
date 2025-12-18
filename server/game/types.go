package game

import (
	"github.com/gorilla/websocket"
)

type GameState string

const (
	StateWaiting     GameState = "WAITING"
	StateGuessing    GameState = "GUESSING"
	StateRoundResult GameState = "ROUND_RESULT"
	StateGameOver    GameState = "GAME_OVER"
)

type Item struct {
	Name     string  `json:"name"`
	ImageURL string  `json:"imageUrl"`
	Price    float64 `json:"price"`
}

type Player struct {
	Name         string          `json:"name"`
	Score        int             `json:"score"`
	CurrentGuess float64         `json:"currentGuess"`
	Conn         *websocket.Conn `json:"-"`
	HasGuessed   bool            `json:"hasGuessed"`
	IsHost       bool            `json:"isHost"`
}

type Room struct {
	ID           string             `json:"id"`
	Players      map[string]*Player `json:"players"`
	Items        []Item             `json:"items"`
	CurrentItem  int                `json:"currentItem"`
	State        GameState          `json:"state"`
	TimeLeft     int                `json:"timeLeft"`
	Broadcast    chan []byte        `json:"-"`
	Register     chan *Player       `json:"-"`
	Unregister   chan *Player       `json:"-"`
	ProcessGuess chan GuessMessage  `json:"-"`
	StartGame    chan *Player       `json:"-"`
	ResetGame    chan *Player       `json:"-"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type GuessMessage struct {
	PlayerName string  `json:"playerName"`
	Guess      float64 `json:"guess"`
}

type JoinMessage struct {
	Name     string `json:"name"`
	RoomCode string `json:"roomCode"`
}
