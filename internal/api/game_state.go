package api

import (
	"fmt"
	"time"

	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

// GameStateResponse represents the complete game state
type GameStateResponse struct {
	Game         GameInfo         `json:"game"`
	CurrentRound *RoundInfo       `json:"currentRound,omitempty"`
	CurrentTurn  *TurnInfo        `json:"currentTurn,omitempty"`
	Players      []PlayerInfo     `json:"players"`
	MyCards      []CardInfo       `json:"myCards,omitempty"`
	InPlayCards  []PlayedCardInfo `json:"inPlayCards,omitempty"`
	RecentLogs   []LogInfo        `json:"recentLogs"`
}

type GameInfo struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	RequesterID string     `json:"requesterId"`
	MaxPlayers  int        `json:"maxPlayers"`
	MinPlayers  int        `json:"minPlayers"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	LoserID     *string    `json:"loserId,omitempty"`
}

type RoundInfo struct {
	ID           string     `json:"id"`
	RoundNumber  int        `json:"roundNumber"`
	Status       string     `json:"status"`
	StartedAt    time.Time  `json:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	LoserID      *string    `json:"loserId,omitempty"`
	DiscardCount int        `json:"discardCount"`
	DiscardPile  []CardInfo `json:"discardPile"`
}

type TurnInfo struct {
	ID               string     `json:"id"`
	TurnNumber       int        `json:"turnNumber"`
	StartPlayerID    string     `json:"startPlayerId"`
	LeadSuit         *string    `json:"leadSuit,omitempty"`
	Status           string     `json:"status"`
	WinnerID         *string    `json:"winnerId,omitempty"`
	CutPlayerID      *string    `json:"cutPlayerId,omitempty"`
	StartedAt        time.Time  `json:"startedAt"`
	CompletedAt      *time.Time `json:"completedAt,omitempty"`
	ExpectedPlayerID *string    `json:"expectedPlayerId,omitempty"`
}

type PlayerInfo struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	IsBot         bool      `json:"isBot"`
	BotDifficulty string    `json:"botDifficulty,omitempty"`
	IsConnected   bool      `json:"isConnected"`
	DonkeyLetters string    `json:"donkeyLetters"`
	JoinOrder     int       `json:"joinOrder"`
	Position      *int      `json:"position,omitempty"` // Position in current round
	CardsInHand   int       `json:"cardsInHand"`
	IsFinished    bool      `json:"isFinished"` // Finished current round
	LastSeenAt    time.Time `json:"lastSeenAt"`
}

type CardInfo struct {
	ID        string `json:"id"`
	Suit      string `json:"suit"`
	Rank      string `json:"rank"`
	Value     int    `json:"value"`
	SortOrder int    `json:"sortOrder"`
	Code      string `json:"code"`
}

type PlayedCardInfo struct {
	ID        string    `json:"id"`
	Card      CardInfo  `json:"card"`
	PlayerID  string    `json:"playerId"`
	PlayOrder int       `json:"playOrder"`
	PlayedAt  time.Time `json:"playedAt"`
}

type LogInfo struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"userId,omitempty"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

// buildGameState builds the complete game state response
func buildGameState(gameID, userID string) (*GameStateResponse, error) {
	// Load game
	var game model.Game
	if err := db.DB.First(&game, "id = ?", gameID).Error; err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}

	// Load game players
	var gamePlayers []model.GamePlayer
	if err := db.DB.Preload("User").Where("game_id = ?", gameID).Order("join_order").Find(&gamePlayers).Error; err != nil {
		return nil, fmt.Errorf("failed to load game players: %w", err)
	}

	// Build player info
	var players []PlayerInfo
	for _, gp := range gamePlayers {
		player := PlayerInfo{
			ID:            gp.UserID,
			Name:          gp.User.Name,
			IsBot:         gp.User.IsBot,
			BotDifficulty: gp.User.BotDifficulty,
			IsConnected:   gp.IsConnected,
			DonkeyLetters: gp.DonkeyLetters,
			JoinOrder:     gp.JoinOrder,
			LastSeenAt:    gp.LastSeenAt,
		}
		players = append(players, player)
	}

	response := &GameStateResponse{
		Game: GameInfo{
			ID:          game.ID,
			Status:      game.Status,
			RequesterID: game.RequesterID,
			MaxPlayers:  game.MaxPlayers,
			MinPlayers:  game.MinPlayers,
			StartedAt:   game.StartedAt,
			CompletedAt: game.CompletedAt,
			LoserID:     game.LoserID,
		},
		Players: players,
	}

	// If game is active, load round and turn info
	if game.Status == "active" {
		// Load current round
		var round model.Round
		if err := db.DB.Where("game_id = ? AND status IN ('dealing', 'active')", gameID).
			Order("round_number DESC").First(&round).Error; err == nil {

			// Count discard pile
			var discardCount int64
			db.DB.Model(&model.Card{}).Where("round_id = ? AND location = 'discard'", round.ID).Count(&discardCount)

			// Load actual discarded cards
			var discardedCards []model.Card
			db.DB.Where("round_id = ? AND location = 'discard'", round.ID).
				Order("updated_at ASC"). // Order by when they were discarded
				Find(&discardedCards)

			// Convert to CardInfo
			var discardPile []CardInfo
			for _, card := range discardedCards {
				discardPile = append(discardPile, CardInfo{
					ID:   card.ID,
					Rank: card.Rank,
					Suit: card.Suit,
				})
			}

			response.CurrentRound = &RoundInfo{
				ID:           round.ID,
				RoundNumber:  round.RoundNumber,
				Status:       round.Status,
				StartedAt:    round.StartedAt,
				CompletedAt:  round.CompletedAt,
				LoserID:      round.LoserID,
				DiscardCount: int(discardCount),
				DiscardPile:  discardPile,
			}

			// Update player info with round data
			var roundPlayers []model.RoundPlayer
			if err := db.DB.Where("round_id = ?", round.ID).Find(&roundPlayers).Error; err == nil {
				roundPlayerMap := make(map[string]model.RoundPlayer)
				for _, rp := range roundPlayers {
					roundPlayerMap[rp.UserID] = rp
				}

				for i := range response.Players {
					if rp, exists := roundPlayerMap[response.Players[i].ID]; exists {
						response.Players[i].Position = &rp.Position
						response.Players[i].CardsInHand = rp.CardsInHand
						response.Players[i].IsFinished = rp.IsFinished
					}
				}
			}

            // Load latest turn for this round (active, cut, or completed). This ensures
            // we show the just-finished turn's in-play cards during the 3s pause,
            // and avoids accidentally showing an older CUT turn.
            var turn model.Turn
            if err := db.DB.Where("round_id = ?", round.ID).
                Order("turn_number DESC").
                Preload("PlayedCards.Card").
                Preload("PlayedCards.Player").
                First(&turn).Error; err == nil {

				// Determine expected player
				expectedPlayerID, _ := getExpectedPlayerIDForTurn(&turn)

				response.CurrentTurn = &TurnInfo{
					ID:               turn.ID,
					TurnNumber:       turn.TurnNumber,
					StartPlayerID:    turn.StartPlayerID,
					LeadSuit:         turn.LeadSuit,
					Status:           turn.Status,
					WinnerID:         turn.WinnerID,
					CutPlayerID:      turn.CutPlayerID,
					StartedAt:        turn.StartedAt,
					CompletedAt:      turn.CompletedAt,
					ExpectedPlayerID: &expectedPlayerID,
				}

				// Build in-play cards
				var inPlayCards []PlayedCardInfo
				for _, pc := range turn.PlayedCards {
					playedCard := PlayedCardInfo{
						ID:        pc.ID,
						PlayerID:  pc.PlayerID,
						PlayOrder: pc.PlayOrder,
						PlayedAt:  pc.PlayedAt,
						Card: CardInfo{
							ID:        pc.Card.ID,
							Suit:      pc.Card.Suit,
							Rank:      pc.Card.Rank,
							Value:     pc.Card.Value,
							SortOrder: pc.Card.SortOrder,
							Code:      pc.Card.CardCode(),
						},
					}
					inPlayCards = append(inPlayCards, playedCard)
				}
				response.InPlayCards = inPlayCards
			}

			// Load user's cards
			var userCards []model.Card
			if err := db.DB.Where("round_id = ? AND owner_id = ? AND location = 'hand'", round.ID, userID).
				Order("sort_order").Find(&userCards).Error; err == nil {

				var myCards []CardInfo
				for _, card := range userCards {
					cardInfo := CardInfo{
						ID:        card.ID,
						Suit:      card.Suit,
						Rank:      card.Rank,
						Value:     card.Value,
						SortOrder: card.SortOrder,
						Code:      card.CardCode(),
					}
					myCards = append(myCards, cardInfo)
				}
				response.MyCards = myCards
			}
		}
	}

	// Load recent logs
	var logs []model.GameSessionLog
	if err := db.DB.Where("game_id = ?", gameID).Order("created_at DESC").Limit(20).Find(&logs).Error; err == nil {
		var recentLogs []LogInfo
		for _, log := range logs {
			logInfo := LogInfo{
				ID:        log.ID,
				UserID:    log.UserID,
				Type:      log.Type,
				Message:   log.Message,
				CreatedAt: log.CreatedAt,
			}
			recentLogs = append(recentLogs, logInfo)
		}
		response.RecentLogs = recentLogs
	}

	return response, nil
}

// buildAdminGameState builds the admin game state with all player cards visible
func buildAdminGameState(gameID string) (*GameStateResponse, error) {
	// Build regular state first
	state, err := buildGameState(gameID, "")
	if err != nil {
		return nil, err
	}

	// If game is active, load all players' cards for admin view
	if state.Game.Status == "active" && state.CurrentRound != nil {
		roundID := state.CurrentRound.ID

		// Load all players' cards
		var allCards []model.Card
		if err := db.DB.Where("round_id = ? AND location = 'hand'", roundID).
			Order("owner_id, sort_order").Find(&allCards).Error; err == nil {

			// Group cards by player
			playerCards := make(map[string][]model.Card)
			for _, card := range allCards {
				if card.OwnerID != nil {
					playerCards[*card.OwnerID] = append(playerCards[*card.OwnerID], card)
				}
			}

			// Add all players' cards to response
			state.MyCards = nil // Clear user-specific cards
			for playerID, cards := range playerCards {
				var cardInfos []CardInfo
				for _, card := range cards {
					cardInfo := CardInfo{
						ID:        card.ID,
						Suit:      card.Suit,
						Rank:      card.Rank,
						Value:     card.Value,
						SortOrder: card.SortOrder,
						Code:      card.CardCode(),
					}
					cardInfos = append(cardInfos, cardInfo)
				}
				
				// Find player and add their cards
				for i := range state.Players {
					if state.Players[i].ID == playerID {
						// Store cards in a separate field for admin view
						// We could extend PlayerInfo to include Cards []CardInfo for admin
						break
					}
				}
			}
		}
	}

	return state, nil
}

// Helper function to determine expected player for a turn (duplicated from enhanced_game.go, should be consolidated)
func getExpectedPlayerIDForTurn(turn *model.Turn) (string, error) {
    // Load all players in this round to maintain seating order
    var allPlayers []model.RoundPlayer
    if err := db.DB.Where("round_id = ?", turn.RoundID).
        Order("position asc").Find(&allPlayers).Error; err != nil {
        return "", fmt.Errorf("failed to load round players: %w", err)
    }

    if len(allPlayers) == 0 {
        return "", fmt.Errorf("no players in round")
    }

    // Track active (not finished)
    active := make(map[string]bool)
    for _, rp := range allPlayers {
        if !rp.IsFinished {
            active[rp.UserID] = true
        }
    }
    if len(active) == 0 {
        return "", fmt.Errorf("no active players in round")
    }

    // Find start player's index among all
    startIdx := -1
    for i, rp := range allPlayers {
        if rp.UserID == turn.StartPlayerID {
            startIdx = i
            break
        }
    }
    if startIdx == -1 {
        return "", fmt.Errorf("start player not found in round")
    }

    // If no cards played yet: expected is the first active player from startIdx
    if len(turn.PlayedCards) == 0 {
        for i := 0; i < len(allPlayers); i++ {
            idx := (startIdx + i) % len(allPlayers)
            if active[allPlayers[idx].UserID] {
                return allPlayers[idx].UserID, nil
            }
        }
        return "", fmt.Errorf("no active player to start turn")
    }

    // Otherwise: next clockwise active player who hasn't played yet
    played := make(map[string]bool)
    for _, pc := range turn.PlayedCards {
        played[pc.PlayerID] = true
    }
    for i := 0; i < len(allPlayers); i++ {
        idx := (startIdx + i) % len(allPlayers)
        pid := allPlayers[idx].UserID
        if active[pid] && !played[pid] {
            return pid, nil
        }
    }

    return "", fmt.Errorf("all players have played")
}
