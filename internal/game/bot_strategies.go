package game

import (
	"encoding/json"
	"math/rand"
	"sort"
	"time"

	"github.com/kairodrad/donkey/internal/model"
)

// EasyBot implements a simple bot strategy that knows the rules but doesn't optimize play
type EasyBot struct {
	UserID string
}

// MediumBot implements a more strategic bot that tries to preserve its interests
type MediumBot struct {
	UserID string
}

// DifficultBot implements an advanced bot that remembers cards and strategizes accordingly
type DifficultBot struct {
	UserID string
	Memory []model.BotMemory
}

// ChooseCard implements BotStrategy for EasyBot
func (b *EasyBot) ChooseCard(playerCards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// Easy bot strategy: 
	// 1. Follow suit if possible, play random card of that suit
	// 2. If can't follow suit, play random card (cut)
	// 3. Sometimes makes suboptimal choices to simulate human error

	validCards := b.getValidCards(playerCards, gameState)
	if len(validCards) == 0 {
		// This shouldn't happen, but fallback to first card
		return playerCards[0]
	}

	// 20% chance of making a suboptimal play
	if rand.Float64() < 0.2 && len(validCards) > 1 {
		// Play a random valid card instead of optimal
		return validCards[rand.Intn(len(validCards))]
	}

	// Basic strategy: try to play lower cards when following suit
	if gameState.CurrentTurn != nil && gameState.CurrentTurn.LeadSuit != nil {
		suitCards := b.getCardsOfSuit(validCards, *gameState.CurrentTurn.LeadSuit)
		if len(suitCards) > 0 {
			// Sort by value and play lowest
			sort.Slice(suitCards, func(i, j int) bool {
				return suitCards[i].Value < suitCards[j].Value
			})
			return suitCards[0]
		}
	}

	// If cutting or leading, play random card
	return validCards[rand.Intn(len(validCards))]
}

func (b *EasyBot) GetDifficulty() string {
	return "easy"
}

// ChooseCard implements BotStrategy for MediumBot
func (b *MediumBot) ChooseCard(playerCards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// Medium bot strategy:
	// 1. Try to avoid playing high cards early
	// 2. When cutting, choose suit strategically
	// 3. Try to get rid of dangerous high cards when safe
	// 4. Follow suit optimally when possible

	validCards := b.getValidCards(playerCards, gameState)
	if len(validCards) == 0 {
		return playerCards[0]
	}

	// If leading the turn, avoid high cards unless hand is mostly high cards
	if len(gameState.InPlayCards) == 0 {
		return b.chooseLeadCard(validCards)
	}

	// If following suit
	if gameState.CurrentTurn != nil && gameState.CurrentTurn.LeadSuit != nil {
		suitCards := b.getCardsOfSuit(validCards, *gameState.CurrentTurn.LeadSuit)
		if len(suitCards) > 0 {
			return b.chooseFollowCard(suitCards, gameState)
		}
	}

	// If cutting, choose suit and card strategically
	return b.chooseCutCard(validCards, gameState)
}

func (b *MediumBot) GetDifficulty() string {
	return "medium"
}

// ChooseCard implements BotStrategy for DifficultBot
func (b *DifficultBot) ChooseCard(playerCards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// Difficult bot strategy:
	// 1. Remember which players are void in which suits
	// 2. Track what cards other players collected
	// 3. Avoid playing suits where dangerous players are void
	// 4. Use memory to make optimal cuts and leads

	validCards := b.getValidCards(playerCards, gameState)
	if len(validCards) == 0 {
		return playerCards[0]
	}

	// Update memory based on current game state
	b.updateMemory(gameState)

	// If leading, use memory to choose safe suit
	if len(gameState.InPlayCards) == 0 {
		return b.chooseLeadCardWithMemory(validCards, gameState)
	}

	// If following suit
	if gameState.CurrentTurn != nil && gameState.CurrentTurn.LeadSuit != nil {
		suitCards := b.getCardsOfSuit(validCards, *gameState.CurrentTurn.LeadSuit)
		if len(suitCards) > 0 {
			return b.chooseFollowCardWithMemory(suitCards, gameState)
		}
	}

	// If cutting, use memory to choose optimal cut
	return b.chooseCutCardWithMemory(validCards, gameState)
}

func (b *DifficultBot) GetDifficulty() string {
	return "difficult"
}

// Helper methods for all bots

func (b *EasyBot) getValidCards(playerCards []model.Card, gameState model.GameStateSnapshot) []model.Card {
	return getValidCardsForPlay(playerCards, gameState)
}

func (b *MediumBot) getValidCards(playerCards []model.Card, gameState model.GameStateSnapshot) []model.Card {
	return getValidCardsForPlay(playerCards, gameState)
}

func (b *DifficultBot) getValidCards(playerCards []model.Card, gameState model.GameStateSnapshot) []model.Card {
	return getValidCardsForPlay(playerCards, gameState)
}

// getValidCardsForPlay returns cards that can be legally played
func getValidCardsForPlay(playerCards []model.Card, gameState model.GameStateSnapshot) []model.Card {
	if gameState.CurrentTurn == nil || gameState.CurrentTurn.LeadSuit == nil {
		// Can play any card when leading
		return playerCards
	}

	// Must follow suit if possible
	leadSuit := *gameState.CurrentTurn.LeadSuit
	suitCards := []model.Card{}
	
	for _, card := range playerCards {
		if card.Suit == leadSuit {
			suitCards = append(suitCards, card)
		}
	}

	if len(suitCards) > 0 {
		return suitCards
	}

	// If no cards of lead suit, can play any card (cut)
	return playerCards
}

func (b *EasyBot) getCardsOfSuit(cards []model.Card, suit string) []model.Card {
	return getCardsOfSuit(cards, suit)
}

func (b *MediumBot) getCardsOfSuit(cards []model.Card, suit string) []model.Card {
	return getCardsOfSuit(cards, suit)
}

func (b *DifficultBot) getCardsOfSuit(cards []model.Card, suit string) []model.Card {
	return getCardsOfSuit(cards, suit)
}

func getCardsOfSuit(cards []model.Card, suit string) []model.Card {
	var suitCards []model.Card
	for _, card := range cards {
		if card.Suit == suit {
			suitCards = append(suitCards, card)
		}
	}
	return suitCards
}

// Medium bot strategy methods

func (b *MediumBot) chooseLeadCard(cards []model.Card) model.Card {
	// When leading, prefer lower cards, but not too low to avoid easy cuts
	// Target middle-range cards when possible
	
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})

	// Try to find a card in the 6-10 range
	for _, card := range cards {
		if card.Value >= 6 && card.Value <= 10 {
			return card
		}
	}

	// If no middle cards, play lowest
	return cards[0]
}

func (b *MediumBot) chooseFollowCard(suitCards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// When following suit, try to play just high enough to win if possible,
	// or lowest if can't win
	
	highestInPlay := b.getHighestCardInPlay(gameState)
	if highestInPlay == nil {
		// First to play in this suit, play low
		sort.Slice(suitCards, func(i, j int) bool {
			return suitCards[i].Value < suitCards[j].Value
		})
		return suitCards[0]
	}

	// Try to play just higher than current highest
	sort.Slice(suitCards, func(i, j int) bool {
		return suitCards[i].Value < suitCards[j].Value
	})

	for _, card := range suitCards {
		if card.Value > highestInPlay.Value {
			return card // Play the lowest card that beats current highest
		}
	}

	// Can't beat it, play lowest
	return suitCards[0]
}

func (b *MediumBot) chooseCutCard(cards []model.Card, _ model.GameStateSnapshot) model.Card {
	// When cutting, prefer to play medium-value cards
	// Avoid aces and kings unless hand is mostly high cards
	
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})

	// Count high cards (J, Q, K, A)
	highCards := 0
	for _, card := range cards {
		if card.Value >= 11 {
			highCards++
		}
	}

	// If most cards are high, might as well play a high one
	if float64(highCards)/float64(len(cards)) > 0.7 {
		// Play highest card to potentially win
		return cards[len(cards)-1]
	}

	// Otherwise try to play a medium card
	for _, card := range cards {
		if card.Value >= 7 && card.Value <= 10 {
			return card
		}
	}

	// Fallback to lowest
	return cards[0]
}

func (b *MediumBot) getHighestCardInPlay(gameState model.GameStateSnapshot) *model.Card {
	if gameState.CurrentTurn == nil || gameState.CurrentTurn.LeadSuit == nil {
		return nil
	}

	var highest *model.Card
	leadSuit := *gameState.CurrentTurn.LeadSuit

	for _, playedCard := range gameState.InPlayCards {
		if playedCard.Card.Suit == leadSuit {
			if highest == nil || playedCard.Card.Value > highest.Value {
				highest = &playedCard.Card
			}
		}
	}

	return highest
}

// Difficult bot strategy methods with memory

func (b *DifficultBot) updateMemory(gameState model.GameStateSnapshot) {
	// Analyze the current turn to update memory about player behaviors
	
	// Track when players cut (are void in a suit)
	if gameState.CurrentTurn != nil && gameState.CurrentTurn.LeadSuit != nil {
		leadSuit := *gameState.CurrentTurn.LeadSuit
		
		for _, playedCard := range gameState.InPlayCards {
			if playedCard.Card.Suit != leadSuit {
				// This player cut, so they're void in the lead suit
				memory := model.SuitVoidMemory{
					PlayerID:   playedCard.PlayerID,
					Suit:       leadSuit,
					TurnNumber: gameState.CurrentTurn.TurnNumber,
				}
				b.addMemory("suit_void", memory, 1.0, nil)
			}
		}
	}

	// Track when players collect cards (could be useful for avoiding those players)
	if gameState.CurrentTurn != nil && gameState.CurrentTurn.Status == "cut" && gameState.CurrentTurn.WinnerID != nil {
		cardCodes := []string{}
		for _, playedCard := range gameState.InPlayCards {
			cardCodes = append(cardCodes, playedCard.Card.CardCode())
		}
		
		memory := model.CardsCollectedMemory{
			PlayerID:   *gameState.CurrentTurn.WinnerID,
			Cards:      cardCodes,
			FromTurnID: gameState.CurrentTurn.ID,
		}
		b.addMemory("cards_collected", memory, 0.9, nil)
	}
}

func (b *DifficultBot) addMemory(memoryType string, data interface{}, confidence float64, expiresAt *time.Time) {
	jsonData, _ := json.Marshal(data)
	
	memory := model.BotMemory{
		ID:          model.NewID(),
		BotUserID:   b.UserID,
		MemoryType:  memoryType,
		MemoryData:  string(jsonData),
		Confidence:  confidence,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}
	
	b.Memory = append(b.Memory, memory)
}

func (b *DifficultBot) chooseLeadCardWithMemory(cards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// Use memory to avoid leading suits where dangerous players are void
	
	// Get players who are close to losing (many DONKEY letters)
	dangerousPlayers := b.getDangerousPlayers(gameState)
	
	// Find suits that dangerous players are NOT void in
	safeSuits := b.getSafeSuits(dangerousPlayers)
	
	// Prefer cards in safe suits
	for _, suit := range safeSuits {
		suitCards := getCardsOfSuit(cards, suit)
		if len(suitCards) > 0 {
			// Within safe suit, use medium bot strategy
			bot := &MediumBot{UserID: b.UserID}
			return bot.chooseLeadCard(suitCards)
		}
	}

	// No safe suits found, fall back to medium strategy
	bot := &MediumBot{UserID: b.UserID}
	return bot.chooseLeadCard(cards)
}

func (b *DifficultBot) chooseFollowCardWithMemory(suitCards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// Use memory to play optimally when following suit
	// For now, use medium bot strategy but could be enhanced with more memory analysis
	bot := &MediumBot{UserID: b.UserID}
	return bot.chooseFollowCard(suitCards, gameState)
}

func (b *DifficultBot) chooseCutCardWithMemory(cards []model.Card, gameState model.GameStateSnapshot) model.Card {
	// When cutting, consider what we know about other players' collections
	// Try to avoid giving cards to players who are doing well
	
	// For now, use medium bot strategy but could be enhanced
	bot := &MediumBot{UserID: b.UserID}
	return bot.chooseCutCard(cards, gameState)
}

func (b *DifficultBot) getDangerousPlayers(gameState model.GameStateSnapshot) []string {
	var dangerous []string
	
	for playerID, letters := range gameState.DonkeyStatus {
		if len(letters) >= 3 { // Players with 3+ letters are dangerous
			dangerous = append(dangerous, playerID)
		}
	}
	
	return dangerous
}

func (b *DifficultBot) getSafeSuits(dangerousPlayers []string) []string {
	suits := []string{"diamonds", "clubs", "hearts", "spades"}
	var safeSuits []string
	
	for _, suit := range suits {
		isSafe := true
		
		// Check if any dangerous player is void in this suit
		for _, memory := range b.Memory {
			if memory.MemoryType == "suit_void" {
				var voidMemory model.SuitVoidMemory
				json.Unmarshal([]byte(memory.MemoryData), &voidMemory)
				
				if voidMemory.Suit == suit {
					for _, dangerousPlayer := range dangerousPlayers {
						if voidMemory.PlayerID == dangerousPlayer {
							isSafe = false
							break
						}
					}
				}
				
				if !isSafe {
					break
				}
			}
		}
		
		if isSafe {
			safeSuits = append(safeSuits, suit)
		}
	}
	
	return safeSuits
}

// Factory function to create bot strategies
func CreateBotStrategy(difficulty string, userID string) model.BotStrategy {
	switch difficulty {
	case "easy":
		return &EasyBot{UserID: userID}
	case "medium":
		return &MediumBot{UserID: userID}
	case "difficult":
		return &DifficultBot{UserID: userID, Memory: []model.BotMemory{}}
	default:
		return &EasyBot{UserID: userID}
	}
}

// Bot name generator for humor
var botNames = []string{
	"ChuckleBot", "GiggleGear", "WittyWidget", "JestJockey", "PunnyPal",
	"LaughLink", "ComedyCog", "SnickerSpark", "HumorHub", "JokesterJet",
	"QuipQueen", "BanterBot", "ChortleChip", "GuffawGuru", "TitterTech",
	"MirthMachine", "CackleCraft", "WhimsyWire", "DrollDroid", "JollyJack",
}

func GenerateRandomBotName() string {
	return botNames[rand.Intn(len(botNames))]
}