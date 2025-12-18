package game

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"price-is-right-server/config"
)

func NewRoom(id string) *Room {
	return &Room{
		ID:           id,
		Players:      make(map[string]*Player),
		Items:        generateItems(),
		CurrentItem:  0,
		State:        StateWaiting,
		TimeLeft:     config.AppConfig.GuessingTime,
		Broadcast:    make(chan []byte),
		Register:     make(chan *Player),
		Unregister:   make(chan *Player),
		ProcessGuess: make(chan GuessMessage),
		StartGame:    make(chan *Player),
		ResetGame:    make(chan *Player),
	}
}

func generateItems() []Item {
	return []Item{
		{Name: "Vintage Toaster", ImageURL: "https://placehold.co/400x300?text=Toaster", Price: 45.0},
		{Name: "Gaming Chair", ImageURL: "https://placehold.co/400x300?text=Chair", Price: 199.99},
		{Name: "Electric Scooter", ImageURL: "https://placehold.co/400x300?text=Scooter", Price: 450.0},
		{Name: "Smart Watch", ImageURL: "https://placehold.co/400x300?text=Watch", Price: 250.0},
		{Name: "Blender", ImageURL: "https://placehold.co/400x300?text=Blender", Price: 89.99},
	}
}

func (r *Room) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case player := <-r.Register:
			if len(r.Players) == 0 {
				player.IsHost = true
			}
			r.Players[player.Name] = player
			r.broadcastState()

		case player := <-r.Unregister:
			if p, ok := r.Players[player.Name]; ok {
				wasHost := p.IsHost
				delete(r.Players, player.Name)
				
				if wasHost && len(r.Players) > 0 {
					// Assign new host to a random remaining player
					for _, newHost := range r.Players {
						newHost.IsHost = true
						break
					}
				}
				// close(player.Broadcast)
			}
			r.broadcastState()

		case player := <-r.StartGame:
			if player.IsHost && r.State == StateWaiting {
				r.State = StateGuessing
				r.CurrentItem = 0
				r.TimeLeft = config.AppConfig.GuessingTime
				r.resetGuesses()
				r.broadcastState()
			}

		case player := <-r.ResetGame:
			if player.IsHost {
				r.State = StateWaiting
				r.CurrentItem = 0
				r.TimeLeft = config.AppConfig.GuessingTime
				r.resetGuesses()
				r.resetScores()
				r.broadcastState()
			}

		case guess := <-r.ProcessGuess:
			if r.State == StateGuessing {
				if p, ok := r.Players[guess.PlayerName]; ok {
					p.CurrentGuess = guess.Guess
					p.HasGuessed = true
					r.checkAllGuessed()
					r.broadcastState()
				}
			}

		case <-ticker.C:
			if r.State == StateGuessing {
				r.TimeLeft--
				if r.TimeLeft <= 0 {
					r.endRound()
				}
				r.broadcastState()
			} else if r.State == StateRoundResult {
				r.TimeLeft--
				if r.TimeLeft <= 0 {
					r.nextRound()
				}
				r.broadcastState() // Update timer for result screen if we want to show "Next round in X"
			}
		}
	}
}

func (r *Room) resetGuesses() {
	for _, p := range r.Players {
		p.HasGuessed = false
		p.CurrentGuess = 0
	}
}

func (r *Room) resetScores() {
	for _, p := range r.Players {
		p.Score = 0
	}
}

func (r *Room) checkAllGuessed() {
	allGuessed := true
	for _, p := range r.Players {
		if !p.HasGuessed {
			allGuessed = false
			break
		}
	}
	if allGuessed {
		r.endRound()
	}
}

func (r *Room) endRound() {
	r.calculateScores()
	r.State = StateRoundResult
	r.TimeLeft = config.AppConfig.ResultTime
	r.broadcastState()
}

func (r *Room) nextRound() {
	r.CurrentItem++
	if r.CurrentItem >= len(r.Items) {
		r.State = StateGameOver
	} else {
		r.State = StateGuessing
		r.TimeLeft = config.AppConfig.GuessingTime
		r.resetGuesses()
	}
	r.broadcastState()
}

func (r *Room) calculateScores() {
	type playerDiff struct {
		Name string
		Diff float64
	}
	
	var diffs []playerDiff
	targetPrice := r.Items[r.CurrentItem].Price

	for name, p := range r.Players {
		diff := math.Abs(p.CurrentGuess - targetPrice)
		diffs = append(diffs, playerDiff{Name: name, Diff: diff})
	}

	// Sort by difference (ascending)
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Diff < diffs[j].Diff
	})

	// Assign points: 1st = 1, 2nd = 2, etc.
	for i, pd := range diffs {
		if p, ok := r.Players[pd.Name]; ok {
			p.Score += (i + 1)
		}
	}
}

func (r *Room) broadcastState() {
	stateMsg := Message{
		Type:    "STATE_UPDATE",
		Payload: r,
	}
	
	data, err := json.Marshal(stateMsg)
	if err != nil {
		fmt.Println("Error marshaling state:", err)
		return
	}

	for _, p := range r.Players {
		err := p.Conn.WriteMessage(1, data) // 1 is TextMessage
		if err != nil {
			fmt.Println("Error writing to player:", err)
			// Handle disconnection?
		}
	}
}
