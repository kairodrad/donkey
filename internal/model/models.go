package model

// User represents a player in the system.
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:20;not null"`
	Games []Game `gorm:"many2many:game_players;"`
}

// Game represents a card game session.
type Game struct {
	ID          uint `gorm:"primaryKey"`
	RequesterID uint
	Requester   User
	Players     []User `gorm:"many2many:game_players;"`
	HasStarted  bool
	IsComplete  bool
	IsAbandoned bool
	State       GameState `gorm:"foreignKey:GameID"`
}

// GameState holds dynamic information for a game.
type GameState struct {
	GameID  uint         `gorm:"primaryKey"`
	Players []GamePlayer `gorm:"foreignKey:GameID"`
}

// GamePlayer associates a user with a game and their hand.
type GamePlayer struct {
	GameID    uint `gorm:"primaryKey"`
	UserID    uint `gorm:"primaryKey"`
	JoinOrder int
	User      User
	Cards     []GameCard `gorm:"foreignKey:GameID,UserID;references:GameID,UserID"`
}

// GameCard represents a single card assigned to a player in a game.
type GameCard struct {
	ID     uint `gorm:"primaryKey"`
	GameID uint
	UserID uint
	Code   string `gorm:"size:3"`
}
