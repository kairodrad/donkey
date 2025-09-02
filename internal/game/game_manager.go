package game

import (
	"errors"
	"fmt"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

// StatePublisher is a function type for publishing game state updates
type StatePublisher func(gameID string)

// LogPublisher is a function type for publishing log events
type LogPublisher func(gameID, userID, logType, message string, eventData interface{})

// Global state publisher function (set by API package)
var globalStatePublisher StatePublisher

// Global log publisher function (set by API package)
var globalLogPublisher LogPublisher

// SetStatePublisher sets the global state publisher function
func SetStatePublisher(publisher StatePublisher) {
	globalStatePublisher = publisher
}

// SetLogPublisher sets the global log publisher function
func SetLogPublisher(publisher LogPublisher) {
	globalLogPublisher = publisher
}

// publishState publishes game state if publisher is set
func publishState(gameID string) {
	if globalStatePublisher != nil {
		globalStatePublisher(gameID)
	}
}

// GameManager handles the overall game lifecycle
type GameManager struct {
	GameID string
}

// NewGameManager creates a new game manager for the specified game
func NewGameManager(gameID string) *GameManager {
	return &GameManager{GameID: gameID}
}

// StartGame initializes the first round and begins gameplay
func (gm *GameManager) StartGame() error {
	// Load game and validate it can be started
	var game model.Game
	if err := db.DB.Preload("GamePlayers.User").First(&game, "id = ?", gm.GameID).Error; err != nil {
		return fmt.Errorf("game not found: %w", err)
	}

	if game.Status != "waiting" {
		return errors.New("game cannot be started from current status")
	}

	// Update game status and start time
	game.Status = "active"
	now := time.Now()
	game.StartedAt = &now
	
	if err := db.DB.Save(&game).Error; err != nil {
		return fmt.Errorf("failed to update game status: %w", err)
	}

	// Create the first round
	if err := gm.StartNewRound(1); err != nil {
		return fmt.Errorf("failed to start first round: %w", err)
	}

	return nil
}

// StartNewRound creates a new round and deals cards
func (gm *GameManager) StartNewRound(roundNumber int) error {
	// Create round
	round := model.Round{
		ID:          model.NewID(),
		GameID:      gm.GameID,
		RoundNumber: roundNumber,
		Status:      "setup",
		StartedAt:   time.Now(),
	}

	if err := db.DB.Create(&round).Error; err != nil {
		return fmt.Errorf("failed to create round: %w", err)
	}

	// Get active players (not finished the game yet) ordered by join time to ensure consistent positioning
	var gamePlayers []model.GamePlayer
	if err := db.DB.Preload("User").Where("game_id = ? AND donkey_letters != 'DONKEY'", gm.GameID).Order("joined_at ASC").Find(&gamePlayers).Error; err != nil {
		return fmt.Errorf("failed to load players: %w", err)
	}

	// Find requester (game creator) to position them last in turn order
	var game model.Game
	if err := db.DB.First(&game, "id = ?", gm.GameID).Error; err != nil {
		return fmt.Errorf("failed to load game: %w", err)
	}

	// Separate requester from other players and assign positions to match visual top-to-bottom order
	var orderedPlayers []model.GamePlayer
	var requester *model.GamePlayer
	
	for _, gp := range gamePlayers {
		if gp.UserID == game.RequesterID {
			requester = &gp
		} else {
			orderedPlayers = append(orderedPlayers, gp)
		}
	}
	
	// Add requester last to maintain clockwise turn order (opponents top-to-bottom, then current player)
	if requester != nil {
		orderedPlayers = append(orderedPlayers, *requester)
	}

	// Create round players with positions that match visual display order
	for i, gp := range orderedPlayers {
		roundPlayer := model.RoundPlayer{
			RoundID:     round.ID,
			UserID:      gp.UserID,
			User:        gp.User,
			Position:    i,  // Position matches visual order: opponents 0..N-1, requester at N
			IsFinished:  false,
			CardsInHand: 0,
		}
		if err := db.DB.Create(&roundPlayer).Error; err != nil {
			return fmt.Errorf("failed to create round player: %w", err)
		}
	}

	// Create and shuffle deck
	cards := model.CreateStandardDeck(round.ID)
	gm.shuffleCards(cards)

	// Deal cards to players
	if err := gm.dealCardsToPlayers(cards, gamePlayers); err != nil {
		return fmt.Errorf("failed to deal cards: %w", err)
	}

	// Update round status
	round.Status = "dealing"
	if err := db.DB.Save(&round).Error; err != nil {
		return fmt.Errorf("failed to update round status: %w", err)
	}

	// Find player with Ace of Spades to start the first turn
	startPlayerID, err := gm.findPlayerWithAceOfSpades(round.ID)
	if err != nil {
		return fmt.Errorf("failed to find starting player: %w", err)
	}

	// Create first turn
	if err := gm.startFirstTurn(round.ID, startPlayerID); err != nil {
		return fmt.Errorf("failed to start first turn: %w", err)
	}

	// Update round status to active
	round.Status = "active"
	if err := db.DB.Save(&round).Error; err != nil {
		return fmt.Errorf("failed to activate round: %w", err)
	}

	// Log round start
	logMessage := fmt.Sprintf("Round %d started. %d players.", roundNumber, len(gamePlayers))
	if err := gm.logEvent("round_event", logMessage, nil); err != nil {
		return fmt.Errorf("failed to log round start: %w", err)
	}

	return nil
}

// shuffleCards shuffles the deck
func (gm *GameManager) shuffleCards(cards []model.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

// dealCardsToPlayers distributes cards evenly among players
func (gm *GameManager) dealCardsToPlayers(cards []model.Card, gamePlayers []model.GamePlayer) error {
	if len(gamePlayers) == 0 {
		return errors.New("no players to deal to")
	}

	// Start dealing from a random player (as per rules)
	startIdx := rand.Intn(len(gamePlayers))
	
	for i, card := range cards {
		playerIdx := (startIdx + i) % len(gamePlayers)
		playerID := gamePlayers[playerIdx].UserID

		// Update card ownership
		card.Location = "hand"
		card.OwnerID = &playerID

		if err := db.DB.Create(&card).Error; err != nil {
			return fmt.Errorf("failed to create card: %w", err)
		}

		// Update player's card count
		if err := db.DB.Model(&model.RoundPlayer{}).
			Where("round_id = ? AND user_id = ?", card.RoundID, playerID).
			UpdateColumn("cards_in_hand", db.DB.Model(&model.RoundPlayer{}).
				Select("cards_in_hand + 1").
				Where("round_id = ? AND user_id = ?", card.RoundID, playerID)).Error; err != nil {
			return fmt.Errorf("failed to update card count: %w", err)
		}
	}

	return nil
}

// findPlayerWithAceOfSpades finds who has the Ace of Spades to start the round
func (gm *GameManager) findPlayerWithAceOfSpades(roundID string) (string, error) {
	var card model.Card
	if err := db.DB.Where("round_id = ? AND suit = 'spades' AND rank = 'A'", roundID).First(&card).Error; err != nil {
		return "", fmt.Errorf("ace of spades not found: %w", err)
	}

	if card.OwnerID == nil {
		return "", errors.New("ace of spades has no owner")
	}

	return *card.OwnerID, nil
}

// startFirstTurn creates the first turn of the round
func (gm *GameManager) startFirstTurn(roundID, startPlayerID string) error {
	turn := model.Turn{
		ID:            model.NewID(),
		RoundID:       roundID,
		TurnNumber:    1,
		StartPlayerID: startPlayerID,
		Status:        "active",
		StartedAt:     time.Now(),
	}

	if err := db.DB.Create(&turn).Error; err != nil {
		return fmt.Errorf("failed to create turn: %w", err)
	}

	// Log turn start
	logMessage := fmt.Sprintf("Turn 1 started. Player %s has the Ace of Spades.", startPlayerID)
	if err := gm.logEvent("turn_event", logMessage, nil); err != nil {
		return fmt.Errorf("failed to log turn start: %w", err)
	}

	// Publish initial active turn so clients can render expected player
	publishState(gm.GameID)

	// Trigger bot play if the starting player is a bot
	// Start the turn sequence (handles both human and bot plays)
	go func() {
		gm.continueTurnSequence(turn.ID)
	}()

	return nil
}

// PlayCard handles a player playing a card
func (gm *GameManager) PlayCard(userID, cardID string) error {
	// Load current game state
	currentTurn, err := gm.getCurrentTurn()
	if err != nil {
		return fmt.Errorf("failed to get current turn: %w", err)
	}

	if currentTurn.Status != "active" {
		return errors.New("turn is not active")
	}

	// Validate it's the player's turn
	expectedPlayerID, err := gm.getExpectedPlayerID(currentTurn)
	if err != nil {
		return fmt.Errorf("failed to determine expected player: %w", err)
	}

	if userID != expectedPlayerID {
		return errors.New("not your turn")
	}

	// Validate card ownership and legality
	if err := gm.validateCardPlay(userID, cardID, currentTurn); err != nil {
		return fmt.Errorf("invalid card play: %w", err)
	}

	// Execute the card play
	if err := gm.executeCardPlay(userID, cardID, currentTurn); err != nil {
		return fmt.Errorf("failed to execute card play: %w", err)
	}

	// Publish state immediately so frontend can see the move
	publishState(gm.GameID)

	// Check if this card play resulted in a CUT
	// Reload turn to get the updated PlayedCards
	updatedTurn, err := gm.getCurrentTurn()
	if err != nil {
		return fmt.Errorf("failed to reload turn after card play: %w", err)
	}

	isCut := gm.isTurnCut(updatedTurn)
	if isCut {
		// CUT occurred - handle it immediately and don't continue turn sequence
		go func() {
			time.Sleep(3 * time.Second)
			if err := gm.handleCutTurn(updatedTurn); err != nil {
				// silently ignore debug output
			}
		}()
		return nil
	}

	// No CUT - continue turn sequence after 3-second pause
	go func() {
		time.Sleep(3 * time.Second)
		gm.continueTurnSequence(currentTurn.ID)
	}()

	return nil
}

// getCurrentTurn gets the active turn
func (gm *GameManager) getCurrentTurn() (*model.Turn, error) {
	var turn model.Turn
	if err := db.DB.Where("round_id IN (SELECT id FROM rounds WHERE game_id = ?) AND status = 'active'", gm.GameID).
		Preload("PlayedCards.Card").
		Preload("PlayedCards.Player").
		First(&turn).Error; err != nil {
		return nil, fmt.Errorf("no active turn found: %w", err)
	}
	return &turn, nil
}

// getExpectedPlayerID determines whose turn it is
func (gm *GameManager) getExpectedPlayerID(turn *model.Turn) (string, error) {
    // Load all players to preserve circular order
    var allPlayers []model.RoundPlayer
    if err := db.DB.Where("round_id = ?", turn.RoundID).
        Order("position asc").Find(&allPlayers).Error; err != nil {
        return "", fmt.Errorf("failed to load round players: %w", err)
    }

    if len(allPlayers) == 0 {
        return "", errors.New("no players in round")
    }

    // Build set of active (not finished) players
    active := make(map[string]bool)
    for _, rp := range allPlayers {
        if !rp.IsFinished {
            active[rp.UserID] = true
        }
    }
    if len(active) == 0 {
        return "", errors.New("no active players in round")
    }

    // Find start player's index among all players
    startIdx := -1
    for i, rp := range allPlayers {
        if rp.UserID == turn.StartPlayerID {
            startIdx = i
            break
        }
    }
    if startIdx == -1 {
        return "", errors.New("start player not found in round")
    }

    // If no cards played yet, expected player is the first active player starting from startIdx
    if len(turn.PlayedCards) == 0 {
        for i := 0; i < len(allPlayers); i++ {
            idx := (startIdx + i) % len(allPlayers)
            pid := allPlayers[idx].UserID
            if active[pid] {
                return pid, nil
            }
        }
        return "", errors.New("no active player to start turn")
    }

    // Next player in clockwise order who hasn't played yet among active players
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

    return "", errors.New("all players have played")
}

// validateCardPlay checks if the card play is legal
func (gm *GameManager) validateCardPlay(userID, cardID string, turn *model.Turn) error {
	// Check card ownership
	var card model.Card
	if err := db.DB.Where("id = ? AND owner_id = ? AND location = 'hand'", cardID, userID).First(&card).Error; err != nil {
		return errors.New("card not found or not owned by player")
	}

	// Special rule: First turn of first round must be Ace of Spades only
	if turn.TurnNumber == 1 && len(turn.PlayedCards) == 0 {
		// Check if this is first round
		var round model.Round
		if err := db.DB.Where("id = ?", turn.RoundID).First(&round).Error; err != nil {
			return fmt.Errorf("failed to load round: %w", err)
		}
		
		if round.RoundNumber == 1 {
			// First turn of first round: must play Ace of Spades
			if card.Rank != "A" || card.Suit != "spades" {
				return errors.New("first turn of first round must be Ace of Spades")
			}
		}
	}

	// Check suit following rules
	if turn.LeadSuit != nil {
		// Must follow suit if possible
		leadSuit := *turn.LeadSuit
		if card.Suit != leadSuit {
			// Check if player has any cards of the lead suit
			var count int64
			if err := db.DB.Model(&model.Card{}).
				Where("round_id = ? AND owner_id = ? AND suit = ? AND location = 'hand'", 
					card.RoundID, userID, leadSuit).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check suit cards: %w", err)
			}
			
			if count > 0 {
				return errors.New("must follow suit when possible")
			}
			// Player is void in lead suit, can cut with any card
		}
	}

	return nil
}

// executeCardPlay performs the actual card play
func (gm *GameManager) executeCardPlay(userID, cardID string, turn *model.Turn) error {
	// Load the card
	var card model.Card
	if err := db.DB.First(&card, "id = ?", cardID).Error; err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	// Set lead suit if this is the first card
	if turn.LeadSuit == nil {
		turn.LeadSuit = &card.Suit
		if err := db.DB.Save(turn).Error; err != nil {
			return fmt.Errorf("failed to set lead suit: %w", err)
		}
	}

	// Create played card record
	playedCard := model.PlayedCard{
		ID:        model.NewID(),
		TurnID:    turn.ID,
		CardID:    cardID,
		PlayerID:  userID,
		PlayOrder: len(turn.PlayedCards) + 1,
		PlayedAt:  time.Now(),
		Card:      card,
	}

	if err := db.DB.Create(&playedCard).Error; err != nil {
		return fmt.Errorf("failed to create played card: %w", err)
	}

	// Update card location
	card.Location = "in_play"
	if err := db.DB.Save(&card).Error; err != nil {
		return fmt.Errorf("failed to update card location: %w", err)
	}

	// Update player's card count
	if err := db.DB.Model(&model.RoundPlayer{}).
		Where("round_id = ? AND user_id = ?", card.RoundID, userID).
		UpdateColumn("cards_in_hand", db.DB.Model(&model.RoundPlayer{}).
			Select("cards_in_hand - 1").
			Where("round_id = ? AND user_id = ?", card.RoundID, userID)).Error; err != nil {
		return fmt.Errorf("failed to update card count: %w", err)
	}

	// Check if player finished the round (no more cards)
	var cardCount int64
	if err := db.DB.Model(&model.Card{}).
		Where("round_id = ? AND owner_id = ? AND location = 'hand'", card.RoundID, userID).
		Count(&cardCount).Error; err != nil {
		return fmt.Errorf("failed to count remaining cards: %w", err)
	}

	if cardCount == 0 {
		// Player finished the round
		now := time.Now()
		if err := db.DB.Model(&model.RoundPlayer{}).
			Where("round_id = ? AND user_id = ?", card.RoundID, userID).
			Updates(map[string]interface{}{
				"is_finished": true,
				"finished_at": &now,
				"cards_in_hand": 0,
			}).Error; err != nil {
			return fmt.Errorf("failed to mark player finished: %w", err)
		}

		// Log player finishing
		if err := gm.logEvent("round_event", fmt.Sprintf("Player %s finished the round!", userID), nil); err != nil {
			return fmt.Errorf("failed to log player finish: %w", err)
		}
	}

	return nil
}

// continueTurnSequence handles the sequential turn progression with proper timing
func (gm *GameManager) continueTurnSequence(turnID string) {
	for {
		// Load current turn
		var turn model.Turn
		if err := db.DB.Preload("PlayedCards.Card").Preload("PlayedCards.Player").Where("id = ?", turnID).First(&turn).Error; err != nil {
			// silently ignore debug output
			return
		}

		if turn.Status != "active" {
			// Turn is completed or cut, stop the sequence
			return
		}

		// Check if turn was cut - if so, stop sequence immediately
		if gm.isTurnCut(&turn) {
			// Turn was cut, stop the sequence - handleCutTurn should have been called already
			return
		}

		// Check if all active players have played
		var activePlayers []model.RoundPlayer
		if err := db.DB.Where("round_id = ? AND is_finished = false", turn.RoundID).Find(&activePlayers).Error; err != nil {
			// silently ignore debug output
			return
		}

		playedPlayers := make(map[string]bool)
		for _, pc := range turn.PlayedCards {
			playedPlayers[pc.PlayerID] = true
		}

		allPlayed := true
		for _, ap := range activePlayers {
			if !playedPlayers[ap.UserID] {
				allPlayed = false
				break
			}
		}

		if allPlayed {
			// All players have played, complete the turn
			if err := gm.completeTurn(&turn); err != nil {
				// silently ignore debug output
			}
			return
		}

		// Get next player to play
		nextPlayerID, err := gm.getExpectedPlayerID(&turn)
		if err != nil {
			// silently ignore debug output
			return
		}

		// Check if next player is a bot
		var user model.User
		if err := db.DB.First(&user, "id = ?", nextPlayerID).Error; err != nil {
			// silently ignore debug output
			return
		}

		if !user.IsBot {
			// Human player, wait for their input
			return
		}

		// Bot player - make them play immediately
		if err := gm.makeBotPlayCard(nextPlayerID, &turn); err != nil {
			// silently ignore debug output
			return
		}

		// Publish state after bot play
		publishState(gm.GameID)

		// Wait 3 seconds before next iteration
		time.Sleep(3 * time.Second)
	}
}

// makeBotPlayCard makes a bot play a card without async timing
func (gm *GameManager) makeBotPlayCard(botUserID string, turn *model.Turn) error {
	// Get bot user details
	var botUser model.User
	if err := db.DB.First(&botUser, "id = ?", botUserID).Error; err != nil {
		return fmt.Errorf("bot user not found: %w", err)
	}

	// Create bot strategy
	botStrategy := CreateBotStrategy(botUser.BotDifficulty, botUserID)

	// Get bot's cards
	var botCards []model.Card
	if err := db.DB.Where("round_id = ? AND owner_id = ? AND location = 'hand'", turn.RoundID, botUserID).
		Order("sort_order").Find(&botCards).Error; err != nil {
		return fmt.Errorf("failed to load bot cards: %w", err)
	}

	// Build game state snapshot
	gameState, err := gm.buildGameStateSnapshot(turn, botUserID)
	if err != nil {
		return fmt.Errorf("failed to build game state: %w", err)
	}

    // Let bot choose card using strategy
    chosenCard := botStrategy.ChooseCard(botCards, *gameState)

    // Enforce rules for bot plays just like humans
    // 1) If first turn of first round and bot holds Ace of Spades, it MUST play it
    if turn.TurnNumber == 1 && len(turn.PlayedCards) == 0 {
        var round model.Round
        if err := db.DB.Where("id = ?", turn.RoundID).First(&round).Error; err != nil {
            return fmt.Errorf("failed to load round: %w", err)
        }
        if round.RoundNumber == 1 {
            for _, c := range botCards {
                if c.IsAceOfSpades() {
                    chosenCard = c
                    break
                }
            }
        }
    }

    // 2) Validate the chosen card; if invalid, pick the first valid alternative
    if err := gm.validateCardPlay(botUserID, chosenCard.ID, turn); err != nil {
        // Try to find any valid card to play
        found := false
        for _, c := range botCards {
            if gm.validateCardPlay(botUserID, c.ID, turn) == nil {
                chosenCard = c
                found = true
                break
            }
        }
        if !found {
            return fmt.Errorf("bot has no valid card to play: %w", err)
        }
    }

    // Execute the bot's card play
    if err := gm.executeCardPlay(botUserID, chosenCard.ID, turn); err != nil {
        return fmt.Errorf("failed to execute bot card play: %w", err)
    }

	// Log bot play
	logMessage := fmt.Sprintf("Bot %s played %s", botUser.Name, chosenCard.CardCode())
	if err := gm.logEvent("turn_event", logMessage, nil); err != nil {
		return fmt.Errorf("failed to log bot play: %w", err)
	}

	// Check if this bot play resulted in a CUT
	// Reload turn to get the updated PlayedCards
	updatedTurn, err := gm.getCurrentTurn()
	if err != nil {
		return fmt.Errorf("failed to reload turn after bot card play: %w", err)
	}

	isCut := gm.isTurnCut(updatedTurn)
	if isCut {
		// CUT occurred - handle it immediately
		if err := gm.handleCutTurn(updatedTurn); err != nil {
			return fmt.Errorf("failed to handle cut after bot play: %w", err)
		}
	}

	return nil
}


// buildGameStateSnapshot creates a read-only snapshot for bot decision making
func (gm *GameManager) buildGameStateSnapshot(turn *model.Turn, botUserID string) (*model.GameStateSnapshot, error) {
	// Get bot's cards
	var botCards []model.Card
	if err := db.DB.Where("round_id = ? AND owner_id = ? AND location = 'hand'", turn.RoundID, botUserID).
		Order("sort_order").Find(&botCards).Error; err != nil {
		return nil, fmt.Errorf("failed to load bot cards: %w", err)
	}

	// Get all round players
	var roundPlayers []model.RoundPlayer
	if err := db.DB.Preload("User").Where("round_id = ?", turn.RoundID).Find(&roundPlayers).Error; err != nil {
		return nil, fmt.Errorf("failed to load round players: %w", err)
	}

	// Get player hand counts
	playerHands := make(map[string]int)
	for _, rp := range roundPlayers {
		playerHands[rp.UserID] = rp.CardsInHand
	}

	// Get DONKEY status
	var gamePlayers []model.GamePlayer
	if err := db.DB.Where("game_id = ?", gm.GameID).Find(&gamePlayers).Error; err != nil {
		return nil, fmt.Errorf("failed to load game players: %w", err)
	}

	donkeyStatus := make(map[string]string)
	for _, gp := range gamePlayers {
		donkeyStatus[gp.UserID] = gp.DonkeyLetters
	}

	// Count discard pile
	var discardCount int64
	if err := db.DB.Model(&model.Card{}).Where("round_id = ? AND location = 'discard'", turn.RoundID).Count(&discardCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count discard pile: %w", err)
	}

	snapshot := &model.GameStateSnapshot{
		GameID:        gm.GameID,
		RoundID:       turn.RoundID,
		TurnID:        turn.ID,
		CurrentTurn:   turn,
		PlayedCards:   turn.PlayedCards,
		PlayerHands:   playerHands,
		MyCards:       botCards,
		InPlayCards:   turn.PlayedCards,
		DiscardCount:  int(discardCount),
		RoundPlayers:  roundPlayers,
		DonkeyStatus:  donkeyStatus,
	}

	return snapshot, nil
}

// completeTurn handles end-of-turn logic (cut or discard)
func (gm *GameManager) completeTurn(turn *model.Turn) error {
	// Determine if turn was cut or completed normally
	isCut := gm.isTurnCut(turn)

	if isCut {
		return gm.handleCutTurn(turn)
	} else {
		return gm.handleCompleteTurn(turn)
	}
}

// isTurnCut checks if the turn was cut (someone played different suit)
func (gm *GameManager) isTurnCut(turn *model.Turn) bool {
	if turn.LeadSuit == nil {
		return false
	}

	leadSuit := *turn.LeadSuit
	for _, pc := range turn.PlayedCards {
		if pc.Card.Suit != leadSuit {
			return true
		}
	}
	return false
}

// handleCutTurn handles a turn that was cut
func (gm *GameManager) handleCutTurn(turn *model.Turn) error {
	// Find the highest card in the lead suit
	var winnerID string
	var highestCard *model.Card
	leadSuit := *turn.LeadSuit

	for _, pc := range turn.PlayedCards {
		if pc.Card.Suit == leadSuit {
			if highestCard == nil || pc.Card.Value > highestCard.Value {
				highestCard = &pc.Card
				winnerID = pc.PlayerID
			}
		}
	}

	if highestCard == nil {
		return errors.New("no cards in lead suit found")
	}

	// Find who cut (first player to play different suit)
	var cutPlayerID string
	for _, pc := range turn.PlayedCards {
		if pc.Card.Suit != leadSuit {
			cutPlayerID = pc.PlayerID
			break
		}
	}

    // Resolve player names for better logs
    var winnerUser, cutterUser model.User
    _ = db.DB.First(&winnerUser, "id = ?", winnerID).Error
    _ = db.DB.First(&cutterUser, "id = ?", cutPlayerID).Error

    // Update turn status immediately (without transferring cards yet)
    now := time.Now()
    turn.Status = "cut"
    turn.WinnerID = &winnerID
    turn.CutPlayerID = &cutPlayerID
    turn.CompletedAt = &now

    if err := db.DB.Save(turn).Error; err != nil {
        return fmt.Errorf("failed to update turn: %w", err)
    }

    // Log cut and publish state immediately so players can see the CUT notification
    logMessage := fmt.Sprintf("CUT: %s cut; %s collected %d cards.", 
        nonEmptyName(cutterUser.Name, cutPlayerID), nonEmptyName(winnerUser.Name, winnerID), len(turn.PlayedCards))
    eventData := map[string]interface{}{
        "type":           "cut",
        "winnerPlayerId": winnerID,
        "cutPlayerId":    cutPlayerID,
        "cardsCollected": len(turn.PlayedCards),
        "winnerName":     winnerUser.Name,
        "cutPlayerName":  cutterUser.Name,
    }
    if err := gm.logEvent("turn_event", logMessage, eventData); err != nil {
        return fmt.Errorf("failed to log cut: %w", err)
    }
    publishState(gm.GameID)

    // Wait 3 seconds for players to see the CUT notification, then transfer cards
    go func() {
        time.Sleep(3 * time.Second)

        // NOW transfer the cards to winner's hand
        for _, pc := range turn.PlayedCards {
            card := pc.Card
            card.Location = "hand"
            card.OwnerID = &winnerID
            if err := db.DB.Save(&card).Error; err != nil {
                // continue even on error to avoid stalling the game
                continue
            }

            // Update winner's card count (best-effort)
            if err := db.DB.Model(&model.RoundPlayer{}).
                Where("round_id = ? AND user_id = ?", card.RoundID, winnerID).
                UpdateColumn("cards_in_hand", db.DB.Model(&model.RoundPlayer{}).
                    Select("cards_in_hand + 1").
                    Where("round_id = ? AND user_id = ?", card.RoundID, winnerID)).Error; err != nil {
                // continue even on error
                continue
            }
        }
		
        // Publish state again after cards have been transferred
        publishState(gm.GameID)
		
        // Check if round/game should end, otherwise start next turn with cut player
        if roundEnded, _ := gm.checkRoundEndWithResult(turn.RoundID); !roundEnded {
            _ = gm.startNextTurn(turn.RoundID, cutPlayerID)
        }
    }()

	return nil
}

// handleCompleteTurn handles a turn that completed normally (no cut)
func (gm *GameManager) handleCompleteTurn(turn *model.Turn) error {
	// Find highest card and winner
	var winnerID string
	var highestCard *model.Card

	for _, pc := range turn.PlayedCards {
		if highestCard == nil || pc.Card.Value > highestCard.Value {
			highestCard = &pc.Card
			winnerID = pc.PlayerID
		}
	}

    // Update turn with results FIRST (do not move cards yet) so players can see all in-play cards
    now := time.Now()
    turn.Status = "completed"
    turn.WinnerID = &winnerID
    turn.CompletedAt = &now

    if err := db.DB.Save(turn).Error; err != nil {
        return fmt.Errorf("failed to update turn: %w", err)
    }

    // Resolve winner name for friendlier logs
    var winnerUser model.User
    _ = db.DB.First(&winnerUser, "id = ?", winnerID).Error

    // Log completion and publish state so players can see the discard outcome
    logMessage := fmt.Sprintf("Discarded %d cards. %s starts next turn.", 
        len(turn.PlayedCards), nonEmptyName(winnerUser.Name, winnerID))
    eventData := map[string]interface{}{
        "type":          "discard",
        "cardsDiscarded": len(turn.PlayedCards),
        "nextStartId":    winnerID,
        "nextStartName":  winnerUser.Name,
        "winnerPlayerId": winnerID,
        "winnerName":     winnerUser.Name,
    }
    if err := gm.logEvent("turn_event", logMessage, eventData); err != nil {
        return fmt.Errorf("failed to log completion: %w", err)
    }
    publishState(gm.GameID)

    // Wait 3 seconds for players to see all in-play cards, then move them to discard
    go func() {
        time.Sleep(3 * time.Second)

        // Move cards to discard pile after pause
        for _, pc := range turn.PlayedCards {
            card := pc.Card
            card.Location = "discard"
            card.OwnerID = nil
            if err := db.DB.Save(&card).Error; err != nil {
                // silently ignore
                return
            }
        }

        publishState(gm.GameID)

        // Check if round/game should end
        roundEnded, err := gm.checkRoundEndWithResult(turn.RoundID)
        if err != nil {
            // silently ignore debug output
            return
        }

        // If round continues (not ended), start next turn with winner
        if !roundEnded {
            if err := gm.startNextTurn(turn.RoundID, winnerID); err != nil {
                // silently ignore debug output
            }
        }
    }()

	return nil
}

// checkRoundEnd checks if the round should end
func (gm *GameManager) checkRoundEnd(roundID string) error {
	// Count players still in the round
	var activeCount int64
	if err := db.DB.Model(&model.RoundPlayer{}).
		Where("round_id = ? AND is_finished = false", roundID).Count(&activeCount).Error; err != nil {
		return fmt.Errorf("failed to count active players: %w", err)
	}

	if activeCount <= 1 {
		// Round ends
		return gm.endRound(roundID)
	}

	return nil
}

// checkRoundEndWithResult checks if the round should end and returns a boolean result
func (gm *GameManager) checkRoundEndWithResult(roundID string) (bool, error) {
	// Count players still in the round
	var activeCount int64
	if err := db.DB.Model(&model.RoundPlayer{}).
		Where("round_id = ? AND is_finished = false", roundID).Count(&activeCount).Error; err != nil {
		return false, fmt.Errorf("failed to count active players: %w", err)
	}

	if activeCount <= 1 {
		// Round ends
		if err := gm.endRound(roundID); err != nil {
			return true, fmt.Errorf("failed to end round: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// endRound ends the current round and determines loser
func (gm *GameManager) endRound(roundID string) error {
	// Find the last remaining player (loser)
	var loser model.RoundPlayer
	if err := db.DB.Where("round_id = ? AND is_finished = false", roundID).First(&loser).Error; err != nil {
		return fmt.Errorf("failed to find round loser: %w", err)
	}

	// Update round
	var round model.Round
	if err := db.DB.First(&round, "id = ?", roundID).Error; err != nil {
		return fmt.Errorf("failed to load round: %w", err)
	}

	now := time.Now()
	round.Status = "completed"
	round.CompletedAt = &now
	round.LoserID = &loser.UserID

	if err := db.DB.Save(&round).Error; err != nil {
		return fmt.Errorf("failed to update round: %w", err)
	}

	// Add DONKEY letter to loser
	var gamePlayer model.GamePlayer
	if err := db.DB.Where("game_id = ? AND user_id = ?", gm.GameID, loser.UserID).First(&gamePlayer).Error; err != nil {
		return fmt.Errorf("failed to load game player: %w", err)
	}

	gamePlayer.AddDonkeyLetter()
	if err := db.DB.Save(&gamePlayer).Error; err != nil {
		return fmt.Errorf("failed to update player letters: %w", err)
	}

	// Log round end
	logMessage := fmt.Sprintf("Round %d ended. Player %s gets letter '%s' (now: %s)", 
		round.RoundNumber, loser.UserID, string(gamePlayer.DonkeyLetters[len(gamePlayer.DonkeyLetters)-1]), gamePlayer.DonkeyLetters)
	if err := gm.logEvent("round_event", logMessage, nil); err != nil {
		return fmt.Errorf("failed to log round end: %w", err)
	}

	// Publish completed round state before proceeding to game end or next round
	publishState(gm.GameID)

	// Check if game should end
	if gamePlayer.IsDonkey() {
		if err := gm.endGame(loser.UserID); err != nil {
			return err
		}
		publishState(gm.GameID)
		return nil
	}

	// Start next round
	return gm.StartNewRound(round.RoundNumber + 1)
}

// endGame ends the entire game
func (gm *GameManager) endGame(loserID string) error {
	// Update game
	var game model.Game
	if err := db.DB.First(&game, "id = ?", gm.GameID).Error; err != nil {
		return fmt.Errorf("failed to load game: %w", err)
	}

	now := time.Now()
	game.Status = "completed"
	game.CompletedAt = &now
	game.LoserID = &loserID

	if err := db.DB.Save(&game).Error; err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	// Get loser player name
	var loserUser model.User
	db.DB.First(&loserUser, "id = ?", loserID)
	
	// Get all players with their DONKEY letter counts for scoreboard
	var gamePlayers []model.GamePlayer
	if err := db.DB.Preload("User").Where("game_id = ?", gm.GameID).Find(&gamePlayers).Error; err != nil {
		return fmt.Errorf("failed to load game players: %w", err)
	}
	
	// Build scoreboard data
	scoreboard := make([]map[string]interface{}, 0, len(gamePlayers))
	for _, gp := range gamePlayers {
        playerData := map[string]interface{}{
            "playerId":    gp.UserID,
            "playerName":  nonEmptyName(gp.User.Name, gp.UserID),
            "donkeyLetters": gp.DonkeyLetters,
            "isLoser":     gp.UserID == loserID,
            "isBot":       gp.User.IsBot,
        }
		scoreboard = append(scoreboard, playerData)
	}

	// Log game end with structured data
	logMessage := fmt.Sprintf("Game completed! Player %s is the DONKEY!", nonEmptyName(loserUser.Name, loserID))
	eventData := map[string]interface{}{
		"type":       "game_end",
		"loserId":    loserID,
		"loserName":  loserUser.Name,
		"scoreboard": scoreboard,
	}
	if err := gm.logEvent("game_event", logMessage, eventData); err != nil {
		return fmt.Errorf("failed to log game end: %w", err)
	}

	// Publish final game status so clients resume UI from paused state
	publishState(gm.GameID)
	return nil
}

// startNextTurn creates the next turn
func (gm *GameManager) startNextTurn(roundID, startPlayerID string) error {
	// Get current turn number
	var maxTurnNumber int
	if err := db.DB.Model(&model.Turn{}).Where("round_id = ?", roundID).
		Select("COALESCE(MAX(turn_number), 0)").Scan(&maxTurnNumber).Error; err != nil {
		return fmt.Errorf("failed to get max turn number: %w", err)
	}

	// Create next turn
	turn := model.Turn{
		ID:            model.NewID(),
		RoundID:       roundID,
		TurnNumber:    maxTurnNumber + 1,
		StartPlayerID: startPlayerID,
		Status:        "active",
		StartedAt:     time.Now(),
	}

	if err := db.DB.Create(&turn).Error; err != nil {
		return fmt.Errorf("failed to create next turn: %w", err)
	}

	// Publish state so clients see the new active turn and expected player
	publishState(gm.GameID)

	// Trigger bot play if start player is bot
	// Start the turn sequence (handles both human and bot plays)
	go func() {
		gm.continueTurnSequence(turn.ID)
	}()

	// Safety watchdog: ensure bots take the first move if expected to
	go func() {
		// Small delay to allow clients to receive the new turn
		time.Sleep(200 * time.Millisecond)
		var checkTurn model.Turn
		if err := db.DB.Preload("PlayedCards").First(&checkTurn, "id = ?", turn.ID).Error; err != nil {
			return
		}
		if checkTurn.Status != "active" || len(checkTurn.PlayedCards) > 0 {
			return
		}
		// Determine expected player
		nextID, err := gm.getExpectedPlayerID(&checkTurn)
		if err != nil || nextID == "" {
			return
		}
		// Check if expected player is a bot
		var u model.User
		if err := db.DB.First(&u, "id = ?", nextID).Error; err != nil {
			return
		}
		if !u.IsBot {
			return
		}
		// Execute bot play
		_ = gm.makeBotPlayCard(nextID, &checkTurn)
		publishState(gm.GameID)
		// Continue sequence after pause
		time.Sleep(3 * time.Second)
		gm.continueTurnSequence(checkTurn.ID)
	}()

	return nil
}

// logEvent creates a log entry and publishes it via SSE if publisher is set
func (gm *GameManager) logEvent(eventType, message string, eventData interface{}) error {
    log := model.GameSessionLog{
        ID:        model.NewID(),
        GameID:    gm.GameID,
        Type:      eventType,
        Message:   message,
        CreatedAt: time.Now(),
    }

    if eventData != nil {
        if b, err := json.Marshal(eventData); err == nil {
            log.EventData = string(b)
        }
    }

    if err := db.DB.Create(&log).Error; err != nil {
        return err
    }

    // Publish log via SSE if publisher is available
    if globalLogPublisher != nil {
        globalLogPublisher(gm.GameID, "", eventType, message, eventData)
    }

    return nil
}

// nonEmptyName returns fallback to id if name is empty
func nonEmptyName(name, id string) string {
    if name == "" {
        return id
    }
    return name
}

// AddBotPlayer adds a bot player to the game
func (gm *GameManager) AddBotPlayer(difficulty string) (*model.User, error) {
	// Create bot user
	botUser := model.User{
		ID:            model.NewID(),
		Name:          GenerateRandomBotName(),
		IsBot:         true,
		BotDifficulty: difficulty,
		CreatedAt:     time.Now(),
	}

	if err := db.DB.Create(&botUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create bot user: %w", err)
	}

	// Join the game
	var game model.Game
	if err := db.DB.First(&game, "id = ?", gm.GameID).Error; err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}

	if game.Status != "waiting" {
		return nil, errors.New("cannot add bot to active game")
	}

	// Get current player count
	var count int64
	if err := db.DB.Model(&model.GamePlayer{}).Where("game_id = ?", gm.GameID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count players: %w", err)
	}

	if count >= 8 {
		return nil, errors.New("game is full")
	}

	// Add to game
	gamePlayer := model.GamePlayer{
		GameID:        gm.GameID,
		UserID:        botUser.ID,
		JoinOrder:     int(count),
		IsConnected:   true,
		DonkeyLetters: "",
		JoinedAt:      time.Now(),
		LastSeenAt:    time.Now(),
	}

	if err := db.DB.Create(&gamePlayer).Error; err != nil {
		return nil, fmt.Errorf("failed to add bot to game: %w", err)
	}

	// Log bot join
	logMessage := fmt.Sprintf("Bot %s (%s) joined the game", botUser.Name, difficulty)
	if err := gm.logEvent("game_event", logMessage, nil); err != nil {
		return nil, fmt.Errorf("failed to log bot join: %w", err)
	}

	return &botUser, nil
}
