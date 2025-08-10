package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// User represents a player in the system.
type User struct {
	ID    string `gorm:"primaryKey;size:32"`
	Name  string `gorm:"size:20;not null"`
	Games []Game `gorm:"many2many:game_players;"`
}

// Game represents a card game session.
type Game struct {
	ID          string `gorm:"primaryKey;size:32"`
	RequesterID string
	Requester   User
	Players     []User `gorm:"many2many:game_players;"`
	HasStarted  bool
	IsComplete  bool
	IsAbandoned bool
	State       GameState `gorm:"foreignKey:GameID"`
}

// GameState holds dynamic information for a game.
type GameState struct {
	GameID  string       `gorm:"primaryKey;size:32"`
	Players []GamePlayer `gorm:"foreignKey:GameID"`
}

// GamePlayer associates a user with a game and their hand.
type GamePlayer struct {
	GameID      string `gorm:"primaryKey;size:32"`
	UserID      string `gorm:"primaryKey;size:32"`
	JoinOrder   int
	User        User
	IsConnected bool
	Cards       []GameCard `gorm:"foreignKey:GameID,UserID;references:GameID,UserID"`
}

// GameCard represents a single card assigned to a player in a game.
type GameCard struct {
	ID     string `gorm:"primaryKey;size:32"`
	GameID string `gorm:"size:32"`
	UserID string `gorm:"size:32"`
	Code   string `gorm:"size:3"`
}

// GameSessionLog records chat and status updates for a game.
type GameSessionLog struct {
	ID        string `gorm:"primaryKey;size:32"`
	GameID    string `gorm:"size:32;index"`
	UserID    string `gorm:"size:32"`
	Type      string `gorm:"size:10"`
	Message   string `gorm:"not null"`
	CreatedAt time.Time
}

// NewID generates a random hexadecimal ID.
func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
