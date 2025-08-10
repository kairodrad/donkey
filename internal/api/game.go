package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/donkey/internal/db"
	"github.com/example/donkey/internal/game"
	"github.com/example/donkey/internal/model"
)

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Name string `json:"name"`
}

// RegisterHandler registers a new user.
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}
	user := model.User{ID: model.NewID(), Name: req.Name}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name})
}

// StartGameRequest represents body to start a game.
type StartGameRequest struct {
	RequesterID string `json:"requesterId"`
}

// StartGameHandler creates a new game.
func StartGameHandler(c *gin.Context) {
	var req StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RequesterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester"})
		return
	}
	gameModel := model.Game{ID: model.NewID(), RequesterID: req.RequesterID}
	if err := db.DB.Create(&gameModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// create initial state and add requester
	db.DB.Create(&model.GameState{GameID: gameModel.ID})
	gp := model.GamePlayer{GameID: gameModel.ID, UserID: req.RequesterID, JoinOrder: 0}
	db.DB.Create(&gp)
	c.JSON(http.StatusOK, gin.H{"gameId": gameModel.ID})
}

// JoinGameRequest represents join body.
type JoinGameRequest struct {
	GameID string `json:"gameId"`
	UserID string `json:"userId"`
}

// JoinGameHandler adds a user to a game.
func JoinGameHandler(c *gin.Context) {
	var req JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.GameID == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid join"})
		return
	}
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	if gameModel.HasStarted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot join"})
		return
	}
	var count int64
	db.DB.Model(&model.GamePlayer{}).Where("game_id = ?", gameModel.ID).Count(&count)
	if count >= 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot join"})
		return
	}
	var existing model.GamePlayer
	if err := db.DB.Where("game_id = ? AND user_id = ?", gameModel.ID, req.UserID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"gameId": gameModel.ID})
		return
	}
	gp := model.GamePlayer{GameID: gameModel.ID, UserID: req.UserID, JoinOrder: int(count)}
	db.DB.Create(&gp)
	// auto finalize if 8 players
	if count+1 >= 8 {
		gameModel.HasStarted = true
		db.DB.Save(&gameModel)
		game.DealCards(&gameModel)
	}
	c.JSON(http.StatusOK, gin.H{"gameId": gameModel.ID})
}

// FinalizeRequest finalizes players and deals cards.
type FinalizeRequest struct {
	GameID string `json:"gameId"`
	UserID string `json:"userId"`
}

// FinalizeHandler finalizes the game and deals cards.
func FinalizeHandler(c *gin.Context) {
	var req FinalizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	if gameModel.HasStarted {
		c.JSON(http.StatusOK, gin.H{"gameId": gameModel.ID})
		return
	}
	if gameModel.RequesterID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only requester can finalize"})
		return
	}
	gameModel.HasStarted = true
	db.DB.Save(&gameModel)
	if err := game.DealCards(&gameModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"gameId": gameModel.ID})
}

// GameStateHandler returns the state of a game for a user.
func GameStateHandler(c *gin.Context) {
	gameID := c.Query("gameId")
	userID := c.Query("userId")
	if gameID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing params"})
		return
	}
	state, err := game.BuildState(gameID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, state)
}
