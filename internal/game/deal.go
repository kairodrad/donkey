package game

import (
	"sort"

	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

// DealCards shuffles the deck and deals all cards to the players in the game.
func DealCards(game *model.Game) error {
	deck := Deck()
	Shuffle(deck)

	var players []model.GamePlayer
	if err := db.DB.Where("game_id = ?", game.ID).Order("join_order asc").Find(&players).Error; err != nil {
		return err
	}
	if len(players) == 0 {
		return nil
	}

	ids := make([]string, len(players))
	requesterIdx := 0
	for i, p := range players {
		ids[i] = p.UserID
		if p.UserID == game.RequesterID {
			requesterIdx = i
		}
	}

	idx := (requesterIdx + 1) % len(ids)
	for _, card := range deck {
		gc := model.GameCard{ID: model.NewID(), GameID: game.ID, UserID: ids[idx], Code: string(card)}
		if err := db.DB.Create(&gc).Error; err != nil {
			return err
		}
		idx = (idx + 1) % len(ids)
	}
	return nil
}

// PlayerState represents what the UI needs to render for a player.
type PlayerState struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Cards     []string `json:"cards,omitempty"`
	CardCount int      `json:"cardCount,omitempty"`
}

// StateResponse describes the game state returned to the client.
type StateResponse struct {
	GameID      string        `json:"gameId"`
	RequesterID string        `json:"requesterId"`
	Players     []PlayerState `json:"players"`
	HasStarted  bool          `json:"hasStarted"`
}

// BuildState builds a StateResponse for a given game and user.
func BuildState(gameID, userID string) (StateResponse, error) {
	var game model.Game
	if err := db.DB.Preload("State.Players.User").Preload("State.Players.Cards").First(&game, "id = ?", gameID).Error; err != nil {
		return StateResponse{}, err
	}

	resp := StateResponse{GameID: game.ID, RequesterID: game.RequesterID, HasStarted: game.HasStarted}

	players := game.State.Players
	sort.Slice(players, func(i, j int) bool { return players[i].JoinOrder < players[j].JoinOrder })

	for _, p := range players {
		ps := PlayerState{ID: p.UserID, Name: p.User.Name}
		if game.HasStarted {
			if p.UserID == userID {
				for _, c := range p.Cards {
					ps.Cards = append(ps.Cards, c.Code)
				}
			} else {
				ps.CardCount = len(p.Cards)
			}
		}
		resp.Players = append(resp.Players, ps)
	}
	return resp, nil
}
