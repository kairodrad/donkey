package game

import (
	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

// DealCards is deprecated - use GameManager.StartGame() instead
func DealCards(game *model.Game) error {
	return nil // No-op for backward compatibility
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
	Status      string        `json:"status"`
}

// BuildState builds a StateResponse for a given game and user.
func BuildState(gameID, userID string) (StateResponse, error) {
	var game model.Game
	if err := db.DB.First(&game, "id = ?", gameID).Error; err != nil {
		return StateResponse{}, err
	}
	var players []model.GamePlayer
	if err := db.DB.Preload("User").Where("game_id = ?", gameID).Order("join_order asc").Find(&players).Error; err != nil {
		return StateResponse{}, err
	}

	resp := StateResponse{GameID: game.ID, RequesterID: game.RequesterID, Status: game.Status}

	// Get current round if game is active
	var currentRoundID string
	if game.Status == "active" {
		var round model.Round
		if err := db.DB.Where("game_id = ? AND status IN ('dealing', 'active')", gameID).
			Order("round_number DESC").First(&round).Error; err == nil {
			currentRoundID = round.ID
		}
	}

	for _, p := range players {
		ps := PlayerState{ID: p.UserID, Name: p.User.Name}
		
		if game.Status == "active" && currentRoundID != "" {
			if p.UserID == userID {
				// Get user's cards from current round
				var userCards []model.Card
				if err := db.DB.Where("round_id = ? AND owner_id = ? AND location = 'hand'", currentRoundID, userID).
					Order("sort_order").Find(&userCards).Error; err == nil {
					for _, c := range userCards {
						ps.Cards = append(ps.Cards, c.CardCode())
					}
				}
			} else {
				// Count other player's cards
				var cardCount int64
				if err := db.DB.Model(&model.Card{}).
					Where("round_id = ? AND owner_id = ? AND location = 'hand'", currentRoundID, p.UserID).
					Count(&cardCount).Error; err == nil {
					ps.CardCount = int(cardCount)
				}
			}
		}
		resp.Players = append(resp.Players, ps)
	}
	return resp, nil
}

// BuildAdminState returns full state exposing all player cards.
func BuildAdminState(gameID string) (StateResponse, error) {
	var game model.Game
	if err := db.DB.First(&game, "id = ?", gameID).Error; err != nil {
		return StateResponse{}, err
	}
	var players []model.GamePlayer
	if err := db.DB.Preload("User").Where("game_id = ?", gameID).Order("join_order asc").Find(&players).Error; err != nil {
		return StateResponse{}, err
	}

	resp := StateResponse{GameID: game.ID, RequesterID: game.RequesterID, Status: game.Status}

	// Get current round if game is active
	var currentRoundID string
	if game.Status == "active" {
		var round model.Round
		if err := db.DB.Where("game_id = ? AND status IN ('dealing', 'active')", gameID).
			Order("round_number DESC").First(&round).Error; err == nil {
			currentRoundID = round.ID
		}
	}

	for _, p := range players {
		ps := PlayerState{ID: p.UserID, Name: p.User.Name}
		
		if game.Status == "active" && currentRoundID != "" {
			// Get all player's cards for admin view
			var playerCards []model.Card
			if err := db.DB.Where("round_id = ? AND owner_id = ? AND location = 'hand'", currentRoundID, p.UserID).
				Order("sort_order").Find(&playerCards).Error; err == nil {
				for _, c := range playerCards {
					ps.Cards = append(ps.Cards, c.CardCode())
				}
				if len(ps.Cards) == 0 {
					ps.CardCount = len(playerCards)
				}
			}
		}
		resp.Players = append(resp.Players, ps)
	}
	return resp, nil
}
