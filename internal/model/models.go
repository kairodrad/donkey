package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Enhanced data model for the Donkey card game
// Supports Games -> Rounds -> Turns architecture with bot players and complex state management

// User represents a player in the system (human or bot)
type User struct {
	ID          string    `gorm:"primaryKey;size:32" json:"id"`
	Name        string    `gorm:"size:20;not null" json:"name"`
	IsBot       bool      `gorm:"default:false" json:"isBot"`
	BotDifficulty string  `gorm:"size:10" json:"botDifficulty,omitempty"` // "easy", "medium", "difficult"
	CreatedAt   time.Time `json:"createdAt"`
	Games       []Game    `gorm:"many2many:game_players;" json:"games"`
}

// Game represents the overall game session that can contain multiple rounds
type Game struct {
	ID           string    `gorm:"primaryKey;size:32" json:"id"`
	RequesterID  string    `json:"requesterId"`
	Requester    User      `json:"-"`
	Players      []User    `gorm:"many2many:game_players;" json:"-"`
	Status       string    `gorm:"size:20;default:'waiting'" json:"status"` // "waiting", "active", "completed", "abandoned", "paused"
	MaxPlayers   int       `gorm:"default:8" json:"maxPlayers"`
	MinPlayers   int       `gorm:"default:2" json:"minPlayers"`
	CreatedAt    time.Time `json:"createdAt"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	LoserID      *string   `json:"loserId,omitempty"` // Final DONKEY loser
	
	// Relationships
	Rounds       []Round   `gorm:"foreignKey:GameID" json:"rounds"`
	GamePlayers  []GamePlayer `gorm:"foreignKey:GameID" json:"gamePlayers"`
	SessionLogs  []GameSessionLog `gorm:"foreignKey:GameID" json:"sessionLogs"`
}

// GamePlayer associates a user with a game and tracks their overall game progress
type GamePlayer struct {
	GameID       string    `gorm:"primaryKey;size:32"`
	UserID       string    `gorm:"primaryKey;size:32"`
	User         User      `json:"user"`
	JoinOrder    int       `json:"joinOrder"`
	IsConnected  bool      `gorm:"default:true" json:"isConnected"`
	DonkeyLetters string   `gorm:"size:6;default:''" json:"donkeyLetters"` // "D", "DO", "DON", "DONK", "DONKE", "DONKEY"
	JoinedAt     time.Time `json:"joinedAt"`
	LastSeenAt   time.Time `json:"lastSeenAt"`
}

// Round represents a single round within a game
type Round struct {
	ID          string    `gorm:"primaryKey;size:32" json:"id"`
	GameID      string    `gorm:"size:32;index" json:"gameId"`
	RoundNumber int       `json:"roundNumber"`
	Status      string    `gorm:"size:20;default:'setup'" json:"status"` // "setup", "dealing", "active", "completed"
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	LoserID     *string   `json:"loserId,omitempty"` // Player who lost this round
	
	// Relationships
	Turns       []Turn    `gorm:"foreignKey:RoundID" json:"turns"`
	RoundPlayers []RoundPlayer `gorm:"foreignKey:RoundID" json:"roundPlayers"`
	Cards       []Card    `gorm:"foreignKey:RoundID" json:"cards"`
}

// RoundPlayer tracks a player's status within a specific round
type RoundPlayer struct {
	RoundID      string    `gorm:"primaryKey;size:32"`
	UserID       string    `gorm:"primaryKey;size:32"`
	User         User      `json:"user"`
	Position     int       `json:"position"` // Seating position (0-7)
	IsFinished   bool      `gorm:"default:false" json:"isFinished"`
	FinishedAt   *time.Time `json:"finishedAt,omitempty"`
	CardsInHand  int       `gorm:"default:0" json:"cardsInHand"`
}

// Turn represents a single turn sequence where players play cards
type Turn struct {
	ID          string    `gorm:"primaryKey;size:32" json:"id"`
	RoundID     string    `gorm:"size:32;index" json:"roundId"`
	TurnNumber  int       `json:"turnNumber"`
	StartPlayerID string  `json:"startPlayerId"`
	LeadSuit    *string   `gorm:"size:10" json:"leadSuit,omitempty"` // "diamonds", "clubs", "hearts", "spades"
	Status      string    `gorm:"size:20;default:'active'" json:"status"` // "active", "cut", "completed"
	WinnerID    *string   `json:"winnerId,omitempty"` // Player who won the turn (highest card)
	CutPlayerID *string   `json:"cutPlayerId,omitempty"` // Player who cut the suit
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	
	// Relationships
	PlayedCards []PlayedCard `gorm:"foreignKey:TurnID" json:"playedCards"`
}

// Card represents a playing card with its current location
type Card struct {
	ID        string    `gorm:"primaryKey;size:32" json:"id"`
	RoundID   string    `gorm:"size:32;index" json:"roundId"`
	Suit      string    `gorm:"size:10;not null" json:"suit"` // "diamonds", "clubs", "hearts", "spades"
	Rank      string    `gorm:"size:5;not null" json:"rank"`  // "2", "3", ..., "10", "J", "Q", "K", "A"
	Value     int       `json:"value"` // Numeric value for comparison (2=2, J=11, Q=12, K=13, A=14)
	Location  string    `gorm:"size:20;default:'deck'" json:"location"` // "deck", "hand", "in_play", "discard"
	OwnerID   *string   `gorm:"size:32" json:"ownerId,omitempty"` // Current owner (if in hand)
	SortOrder int       `json:"sortOrder"` // For consistent card ordering
}

// PlayedCard represents a card played during a turn
type PlayedCard struct {
	ID        string    `gorm:"primaryKey;size:32" json:"id"`
	TurnID    string    `gorm:"size:32;index" json:"turnId"`
	CardID    string    `gorm:"size:32" json:"cardId"`
	PlayerID  string    `gorm:"size:32" json:"playerId"`
	PlayOrder int       `json:"playOrder"` // Order in which card was played in this turn
	PlayedAt  time.Time `json:"playedAt"`
	
	// Relationships
	Card      Card      `json:"card"`
	Player    User      `json:"player"`
}

// BotMemory stores bot's memory of game events for strategic play
type BotMemory struct {
	ID          string    `gorm:"primaryKey;size:32" json:"id"`
	BotUserID   string    `gorm:"size:32;index" json:"botUserId"`
	GameID      string    `gorm:"size:32;index" json:"gameId"`
	MemoryType  string    `gorm:"size:20" json:"memoryType"` // "player_cut", "cards_collected", "suit_void"
	MemoryData  string    `gorm:"type:text" json:"memoryData"` // JSON data for flexible storage
	Confidence  float64   `gorm:"default:1.0" json:"confidence"` // How confident the bot is about this memory
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"` // Some memories expire
}

// GameSessionLog records chat and status updates for a game (enhanced)
type GameSessionLog struct {
	ID        string    `gorm:"primaryKey;size:32" json:"id"`
	GameID    string    `gorm:"size:32;index" json:"gameId"`
	UserID    *string   `gorm:"size:32" json:"userId,omitempty"` // Null for system messages
	Type      string    `gorm:"size:20" json:"type"` // "chat", "system", "game_event", "round_event", "turn_event"
	Message   string    `gorm:"not null" json:"message"`
	EventData string    `gorm:"type:text" json:"eventData,omitempty"` // JSON data for structured events
	CreatedAt time.Time `json:"createdAt"`
}

// GameSettings holds configurable game parameters
type GameSettings struct {
	GameID              string `gorm:"primaryKey;size:32" json:"gameId"`
	AutoStartAt8Players bool   `gorm:"default:true" json:"autoStartAt8Players"`
	AllowBots           bool   `gorm:"default:true" json:"allowBots"`
	MaxBots             int    `gorm:"default:6" json:"maxBots"`
	TurnTimeoutSeconds  int    `gorm:"default:30" json:"turnTimeoutSeconds"`
	PauseOnDisconnect   bool   `gorm:"default:true" json:"pauseOnDisconnect"`
}

// Helper methods and types for game logic

// CardCode returns the traditional card code (e.g., "AS", "KH", "2D")
func (c *Card) CardCode() string {
	suitMap := map[string]string{
		"spades":   "S",
		"hearts":   "H", 
		"diamonds": "D",
		"clubs":    "C",
	}
	return c.Rank + suitMap[c.Suit]
}

// IsAceOfSpades checks if this card is the Ace of Spades (starting card)
func (c *Card) IsAceOfSpades() bool {
	return c.Suit == "spades" && c.Rank == "A"
}

// GetDonkeyProgress returns the number of letters accumulated (0-6)
func (gp *GamePlayer) GetDonkeyProgress() int {
	return len(gp.DonkeyLetters)
}

// IsDonkey checks if player has accumulated all letters
func (gp *GamePlayer) IsDonkey() bool {
	return gp.DonkeyLetters == "DONKEY"
}

// AddDonkeyLetter adds the next letter in the sequence
func (gp *GamePlayer) AddDonkeyLetter() {
	letters := "DONKEY"
	if len(gp.DonkeyLetters) < len(letters) {
		gp.DonkeyLetters += string(letters[len(gp.DonkeyLetters)])
	}
}

// BotStrategy defines the interface for bot behavior
type BotStrategy interface {
	ChooseCard(playerCards []Card, gameState GameStateSnapshot) Card
	GetDifficulty() string
}

// GameStateSnapshot provides read-only game state for bot decision making
type GameStateSnapshot struct {
	GameID        string                 `json:"gameId"`
	RoundID       string                 `json:"roundId"`
	TurnID        string                 `json:"turnId"`
	CurrentTurn   *Turn                  `json:"currentTurn"`
	PlayedCards   []PlayedCard          `json:"playedCards"`
	PlayerHands   map[string]int        `json:"playerHands"` // PlayerID -> card count
	MyCards       []Card                `json:"myCards"`
	InPlayCards   []PlayedCard          `json:"inPlayCards"`
	DiscardCount  int                   `json:"discardCount"`
	RoundPlayers  []RoundPlayer         `json:"roundPlayers"`
	DonkeyStatus  map[string]string     `json:"donkeyStatus"` // PlayerID -> letters
}

// BotMemoryData structures for different memory types
type PlayerCutMemory struct {
	PlayerID    string `json:"playerId"`
	Suit        string `json:"suit"`
	TurnNumber  int    `json:"turnNumber"`
}

type CardsCollectedMemory struct {
	PlayerID    string   `json:"playerId"`
	Cards       []string `json:"cards"` // Card codes
	FromTurnID  string   `json:"fromTurnId"`
}

type SuitVoidMemory struct {
	PlayerID   string `json:"playerId"`
	Suit       string `json:"suit"`
	TurnNumber int    `json:"turnNumber"`
}

// NewID generates a random hexadecimal ID.
func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Helper function to create standard deck
func CreateStandardDeck(roundID string) []Card {
	suits := []string{"diamonds", "clubs", "hearts", "spades"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	values := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
		"J": 11, "Q": 12, "K": 13, "A": 14,
	}
	
	var cards []Card
	sortOrder := 0
	
	for _, suit := range suits {
		for _, rank := range ranks {
			card := Card{
				ID:        NewID(),
				RoundID:   roundID,
				Suit:      suit,
				Rank:      rank,
				Value:     values[rank],
				Location:  "deck",
				SortOrder: sortOrder,
			}
			cards = append(cards, card)
			sortOrder++
		}
	}
	
	return cards
}