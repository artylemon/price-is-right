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
        return []Item{
                {Name: "Dell Latitude 5420 Business Laptop, 14-Inch FHD(1920x1080) Display, Intel Core i7-1165G7, 32GB RAM, 512 GB SSD, QWERTY Keyboard, Windows 11 Pro", ImageURL: "https://m.media-amazon.com/images/I/61AWVWW92nL._AC_SL1500_.jpg", Price: 399.0},
                {Name: "Visual Studio Enterprise Standard (per user per year)", ImageURL: "https://logos-world.net/wp-content/uploads/2025/04/Visual-Studio-Logo.png", Price: 1919.95},
                {Name: "SQL Server Management Studio", ImageURL: "https://d1jnx9ba8s6j9r.cloudfront.net/blog/wp-content/uploads/2019/10/logo.png", Price: 0.0},
                {Name: "Lucid Chart (per user per month)", ImageURL: "https://logovtor.com/wp-content/uploads/2021/09/lucidchart-logo-vector.png", Price: 10.0},
                {Name: "Postman Professional (per user per month)", ImageURL: "https://blog.postman.com/wp-content/uploads/2015/08/postman-logo-drawing-board-825x510.png", Price: 21.67},
                {Name: "Dell Pro Thunderbolt 4 Dock - WD25TB4 (with VAT)", ImageURL: "https://i.dell.com/is/image/DellContent/content/dam/ss2/product-images/dell-client-products/peripherals/docks/wd25tb4/media-gallery/dock-station-wd25tb4-black-gallery-1.psd?fmt=png-alpha&pscan=auto&scl=1&wid=3653&hei=1444&qlt=100,1&resMode=sharp2&size=3653,1444&chrss=full&imwidth=5000", Price: 268.92},
                {Name: "Dell Pro 24 Adjustable Stand Monitor - E2425HSM", ImageURL: "https://i.dell.com/is/image/DellContent/content/dam/ss2/product-images/dell-client-products/peripherals/monitors/e-series/e2425hsm/media-gallery/monitor-dell-pro-e2425hsm-bk-gallery-1.psd?fmt=png-alpha&pscan=auto&scl=1&wid=4513&hei=4178&qlt=100,1&resMode=sharp2&size=4513,4178&chrss=full&imwidth=5000", Price: 97.20},
                {Name: "Markus Office Chair with Armrests", ImageURL: "https://www.ikea.com/gb/en/images/products/markus-office-chair-vissle-dark-grey__0724714_pe734597_s5.jpg?f=xl", Price: 150.0},
        }
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
