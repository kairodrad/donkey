package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/game"
	"github.com/kairodrad/donkey/internal/model"
)

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
    Name string `json:"name"`
}

// RegisterHandler registers a new user.
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}
    req.Name = strings.TrimSpace(req.Name)
    if req.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
        return
    }
    // Capitalize first letter of each word in the provided name
    req.Name = capitalizeName(req.Name)
    user := model.User{
        ID:        model.NewID(),
        Name:      req.Name,
        IsBot:     false,
        CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name})
}

// capitalizeName converts "john doe" -> "John Doe" in a conservative way without
// introducing extra dependencies. It trims internal spacing to single spaces.
func capitalizeName(s string) string {
    parts := strings.Fields(s)
    for i, p := range parts {
        r := []rune(p)
        if len(r) == 0 {
            continue
        }
        // Uppercase first rune, lowercase the remainder
        first := strings.ToUpper(string(r[0]))
        rest := ""
        if len(r) > 1 {
            rest = strings.ToLower(string(r[1:]))
        }
        parts[i] = first + rest
    }
    return strings.Join(parts, " ")
}

// CreateGameRequest represents the game creation request
type CreateGameRequest struct {
	RequesterID string `json:"requesterId"`
	MaxPlayers  int    `json:"maxPlayers,omitempty"`
	MinPlayers  int    `json:"minPlayers,omitempty"`
}

// CreateGameHandler creates a new game
func CreateGameHandler(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RequesterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester"})
		return
	}

	// Validate user exists
	var user model.User
	if err := db.DB.First(&user, "id = ?", req.RequesterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Set defaults
	if req.MaxPlayers == 0 {
		req.MaxPlayers = 8
	}
	if req.MinPlayers == 0 {
		req.MinPlayers = 2
	}

	// Create game
	gameModel := model.Game{
		ID:          model.NewID(),
		RequesterID: req.RequesterID,
		Status:      "waiting",
		MaxPlayers:  req.MaxPlayers,
		MinPlayers:  req.MinPlayers,
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&gameModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create game settings
	settings := model.GameSettings{
		GameID:              gameModel.ID,
		AutoStartAt8Players: true,
		AllowBots:           true,
		MaxBots:             6,
		TurnTimeoutSeconds:  30,
		PauseOnDisconnect:   true,
	}

	if err := db.DB.Create(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add requester as first player
	gamePlayer := model.GamePlayer{
		GameID:        gameModel.ID,
		UserID:        req.RequesterID,
		JoinOrder:     0,
		IsConnected:   true,
		DonkeyLetters: "",
		JoinedAt:      time.Now(),
		LastSeenAt:    time.Now(),
	}

	if err := db.DB.Create(&gamePlayer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log game creation
	logMessage := fmt.Sprintf("Game created by %s", user.Name)
	logAndSend(gameModel.ID, req.RequesterID, "game_event", logMessage)
	publishState(gameModel.ID)

	c.JSON(http.StatusOK, gin.H{
		"gameId":     gameModel.ID,
		"status":     gameModel.Status,
		"maxPlayers": gameModel.MaxPlayers,
		"minPlayers": gameModel.MinPlayers,
	})
}

// JoinGameRequest represents joining a game
type JoinGameRequest struct {
	GameID string `json:"gameId"`
	UserID string `json:"userId"`
}

// JoinGameHandler adds a user to a game
func JoinGameHandler(c *gin.Context) {
	var req JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.GameID == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid join request"})
		return
	}

	// Validate user exists
	var user model.User
	if err := db.DB.First(&user, "id = ?", req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Validate game exists and can be joined
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if gameModel.Status != "waiting" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot join active game"})
		return
	}

	// Check if already in game
	var existing model.GamePlayer
	if err := db.DB.Where("game_id = ? AND user_id = ?", req.GameID, req.UserID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"gameId": req.GameID, "status": "already_joined"})
		return
	}

	// Check player count
	var playerCount int64
	if err := db.DB.Model(&model.GamePlayer{}).Where("game_id = ?", req.GameID).Count(&playerCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count players"})
		return
	}

	if int(playerCount) >= gameModel.MaxPlayers {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game is full"})
		return
	}

	// Add player to game
	gamePlayer := model.GamePlayer{
		GameID:        req.GameID,
		UserID:        req.UserID,
		JoinOrder:     int(playerCount),
		IsConnected:   true,
		DonkeyLetters: "",
		JoinedAt:      time.Now(),
		LastSeenAt:    time.Now(),
	}

	if err := db.DB.Create(&gamePlayer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log player join
	logMessage := fmt.Sprintf("%s joined the game", user.Name)
	logAndSend(req.GameID, req.UserID, "game_event", logMessage)

	// Check if should auto-start
	var settings model.GameSettings
	if err := db.DB.First(&settings, "game_id = ?", req.GameID).Error; err == nil {
		if settings.AutoStartAt8Players && int(playerCount)+1 >= 8 {
			gm := game.NewGameManager(req.GameID)
			if err := gm.StartGame(); err != nil {
				logAndSend(req.GameID, req.UserID, "game_event", "Failed to auto-start game: "+err.Error())
			} else {
				logAndSend(req.GameID, req.UserID, "game_event", "Game auto-started with 8 players!")
			}
		}
	}

	publishState(req.GameID)
	c.JSON(http.StatusOK, gin.H{"gameId": req.GameID, "status": "joined"})
}

// AddBotRequest represents adding a bot to the game
type AddBotRequest struct {
	GameID     string `json:"gameId"`
	UserID     string `json:"userId"`
	Difficulty string `json:"difficulty"` // "easy", "medium", "difficult"
}

// AddBotHandler adds a bot player to the game
func AddBotHandler(c *gin.Context) {
	var req AddBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate requester
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if gameModel.RequesterID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only game creator can add bots"})
		return
	}

	// Validate difficulty
	if req.Difficulty != "easy" && req.Difficulty != "medium" && req.Difficulty != "difficult" {
		req.Difficulty = "easy" // Default
	}

	// Add bot using game manager
	gm := game.NewGameManager(req.GameID)
	botUser, err := gm.AddBotPlayer(req.Difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publishState(req.GameID)
	c.JSON(http.StatusOK, gin.H{
		"botId":      botUser.ID,
		"botName":    botUser.Name,
		"difficulty": botUser.BotDifficulty,
	})
}

// StartGameRequest represents starting a game
type StartGameRequest struct {
	GameID string `json:"gameId"`
	UserID string `json:"userId"`
}

// StartGameHandler starts a game
func StartGameHandler(c *gin.Context) {
	var req StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate requester
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if gameModel.RequesterID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only game creator can start the game"})
		return
	}

	if gameModel.Status != "waiting" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game cannot be started"})
		return
	}

	// Check minimum players
	var playerCount int64
	if err := db.DB.Model(&model.GamePlayer{}).Where("game_id = ?", req.GameID).Count(&playerCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count players"})
		return
	}

	if int(playerCount) < gameModel.MinPlayers {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("need at least %d players", gameModel.MinPlayers)})
		return
	}

	// Start game using game manager
	gm := game.NewGameManager(req.GameID)
	if err := gm.StartGame(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	publishState(req.GameID)
	c.JSON(http.StatusOK, gin.H{"gameId": req.GameID, "status": "started"})
}

// PlayCardRequest represents playing a card
type PlayCardRequest struct {
	GameID string `json:"gameId"`
	UserID string `json:"userId"`
	CardID string `json:"cardId"`
}

// PlayCardHandler handles a player playing a card
func PlayCardHandler(c *gin.Context) {
	var req PlayCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate game is active
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if gameModel.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game is not active"})
		return
	}

	// Validate player is in game
	var gamePlayer model.GamePlayer
	if err := db.DB.Where("game_id = ? AND user_id = ?", req.GameID, req.UserID).First(&gamePlayer).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "player not in game"})
		return
	}

	// Play card using game manager
	gm := game.NewGameManager(req.GameID)
	if err := gm.PlayCard(req.UserID, req.CardID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publishState(req.GameID)
	c.JSON(http.StatusOK, gin.H{"status": "card_played"})
}

// GameStateHandler returns complete game state for a user
func GameStateHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	userID := c.Param("userId")

	if gameID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	state, err := buildGameState(gameID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

// AdminStateHandler returns full game state for debugging
func AdminStateHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing gameId"})
		return
	}

	state, err := buildAdminGameState(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

// GetGameListHandler returns list of games for a user
func GetGameListHandler(c *gin.Context) {
	userID := c.Query("userId")
	status := c.Query("status")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing userId"})
		return
	}

	query := db.DB.Model(&model.Game{}).
		Joins("JOIN game_players ON games.id = game_players.game_id").
		Where("game_players.user_id = ?", userID)

	if status != "" {
		query = query.Where("games.status = ?", status)
	}

	var games []model.Game
	if err := query.Order("games.created_at DESC").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type GameListItem struct {
		ID          string     `json:"id"`
		Status      string     `json:"status"`
		RequesterID string     `json:"requesterId"`
		PlayerCount int        `json:"playerCount"`
		BotCount    int        `json:"botCount"`
		CreatedAt   time.Time  `json:"createdAt"`
		StartedAt   *time.Time `json:"startedAt,omitempty"`
		CompletedAt *time.Time `json:"completedAt,omitempty"`
		LastActivity time.Time `json:"lastActivity"`
	}

	var result []GameListItem
	for _, game := range games {
		// Count players and bots
		var totalPlayers, botPlayers int64
		db.DB.Model(&model.GamePlayer{}).Where("game_id = ?", game.ID).Count(&totalPlayers)
		db.DB.Table("game_players").
			Joins("JOIN users ON game_players.user_id = users.id").
			Where("game_players.game_id = ? AND users.is_bot = true", game.ID).
			Count(&botPlayers)

		// Get last activity
		var lastLog model.GameSessionLog
		lastActivity := game.CreatedAt
		if err := db.DB.Where("game_id = ?", game.ID).Order("created_at DESC").First(&lastLog).Error; err == nil {
			lastActivity = lastLog.CreatedAt
		}

		item := GameListItem{
			ID:           game.ID,
			Status:       game.Status,
			RequesterID:  game.RequesterID,
			PlayerCount:  int(totalPlayers),
			BotCount:     int(botPlayers),
			CreatedAt:    game.CreatedAt,
			StartedAt:    game.StartedAt,
			CompletedAt:  game.CompletedAt,
			LastActivity: lastActivity,
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"games": result})
}

// AbandonGameRequest represents abandoning a game
type AbandonGameRequest struct {
	GameID string `json:"gameId"`
	UserID string `json:"userId"`
}

// AbandonGameHandler abandons a game (only requester can do this)
func AbandonGameHandler(c *gin.Context) {
	var req AbandonGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate requester
	var game model.Game
	if err := db.DB.First(&game, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	if game.RequesterID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only game creator can abandon"})
		return
	}

	if game.Status == "completed" || game.Status == "abandoned" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game already ended"})
		return
	}

	// Update game status
	game.Status = "abandoned"
	now := time.Now()
	game.CompletedAt = &now

	if err := db.DB.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log abandonment
	logAndSend(req.GameID, req.UserID, "game_event", "Game abandoned by creator")
	publishState(req.GameID)

	c.JSON(http.StatusOK, gin.H{"status": "abandoned"})
}
