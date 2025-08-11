package api

import (
	"errors"
	"net/http"
	"strings"

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
	var user model.User
	if err := db.DB.First(&user, "id = ?", req.RequesterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
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
	logAndSend(gameModel.ID, req.RequesterID, "status", user.Name+" created the game")
	publishState(gameModel.ID)
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
	if _, err := joinGame(req.GameID, req.UserID); err != nil {
		if errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"gameId": req.GameID})
}

var errNotFound = errors.New("not found")

func joinGame(gameID, userID string) (*model.Game, error) {
	var user model.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errNotFound
	}
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", gameID).Error; err != nil {
		return nil, errNotFound
	}
	if gameModel.HasStarted {
		return nil, errors.New("cannot join")
	}
	var count int64
	db.DB.Model(&model.GamePlayer{}).Where("game_id = ?", gameID).Count(&count)
	if count >= 8 {
		return nil, errors.New("cannot join")
	}
	var existing model.GamePlayer
	if err := db.DB.Where("game_id = ? AND user_id = ?", gameID, userID).First(&existing).Error; err == nil {
		return &gameModel, nil
	}
	gp := model.GamePlayer{GameID: gameID, UserID: userID, JoinOrder: int(count)}
	db.DB.Create(&gp)
	logAndSend(gameID, userID, "status", user.Name+": joined the game")
	if count+1 >= 8 {
		gameModel.HasStarted = true
		db.DB.Save(&gameModel)
		game.DealCards(&gameModel)
		logAndSend(gameID, userID, "status", "Cards dealt. Begin game.")
	}
	publishState(gameID)
	return &gameModel, nil
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
	logAndSend(gameModel.ID, req.UserID, "status", "Cards dealt. Begin game.")
	publishState(gameModel.ID)
	c.JSON(http.StatusOK, gin.H{"gameId": gameModel.ID})
}

// GameStateHandler returns the state of a game for a user.
//
// @Summary      Get game state
// @Description  Returns current game state tailored for the requesting user
// @Tags         game
// @Produce      json
// @Param        gameId  path  string  true  "Game ID"
// @Param        userId  path  string  true  "User ID"
// @Success      200  {object}  game.StateResponse
// @Router       /api/game/{gameId}/state/{userId} [get]
func GameStateHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	userID := c.Param("userId")
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

// AdminStateHandler returns full state for debugging.
//
// @Summary      Get full game state
// @Description  Returns complete game information with all player cards
// @Tags         admin
// @Produce      json
// @Param        gameId  path  string  true  "Game ID"
// @Success      200  {object}  game.StateResponse
// @Router       /api/admin/game/{gameId}/state [get]
func AdminStateHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing"})
		return
	}
	state, err := game.BuildAdminState(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, state)
}
