package api

import (
	"encoding/json"
	"html"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

type event struct {
	Type string                `json:"type"`
	Log  *model.GameSessionLog `json:"log,omitempty"`
}

type broker struct {
	mu   sync.Mutex
	subs map[string]map[chan event]struct{}
}

var b = broker{subs: make(map[string]map[chan event]struct{})}

func (br *broker) subscribe(gameID string) chan event {
	ch := make(chan event, 8)
	br.mu.Lock()
	if br.subs[gameID] == nil {
		br.subs[gameID] = make(map[chan event]struct{})
	}
	br.subs[gameID][ch] = struct{}{}
	br.mu.Unlock()
	return ch
}

func (br *broker) unsubscribe(gameID string, ch chan event) {
	br.mu.Lock()
	if m, ok := br.subs[gameID]; ok {
		delete(m, ch)
		if len(m) == 0 {
			delete(br.subs, gameID)
		}
	}
	br.mu.Unlock()
	close(ch)
}

func (br *broker) publish(gameID string, ev event) {
	br.mu.Lock()
	m := br.subs[gameID]
	for ch := range m {
		select {
		case ch <- ev:
		default:
		}
	}
	br.mu.Unlock()
}

func logAndSend(gameID, userID, typ, message string) {
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	entry := model.GameSessionLog{ID: model.NewID(), GameID: gameID, UserID: userIDPtr, Type: typ, Message: message, CreatedAt: time.Now()}
	db.DB.Create(&entry)
	b.publish(gameID, event{Type: "log", Log: &entry})
}

func publishState(gameID string) {
	b.publish(gameID, event{Type: "state"})
}

// PublishState is a public wrapper for publishState (used by game manager)
func PublishState(gameID string) {
	publishState(gameID)
}

// PublishLog publishes a structured log event with eventData (used by game manager)
func PublishLog(gameID, userID, logType, message string, eventData interface{}) {
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	
	entry := model.GameSessionLog{
		ID:        model.NewID(), 
		GameID:    gameID, 
		UserID:    userIDPtr, 
		Type:      logType, 
		Message:   message, 
		CreatedAt: time.Now(),
	}
	
	// Add structured eventData if provided
	if eventData != nil {
		if b, err := json.Marshal(eventData); err == nil {
			entry.EventData = string(b)
		}
	}
	
	db.DB.Create(&entry)
	b.publish(gameID, event{Type: "log", Log: &entry})
}

// StreamHandler provides a long-lived stream of events for a game.
//
// @Summary      Stream game updates
// @Description  Streams session and state change events for a game
// @Tags         events
// @Produce      text/event-stream
// @Param        gameId  path  string  true  "Game ID"
// @Param        userId  path  string  true  "User ID"
// @Success      200  {string}  string  "event stream"
// @Router       /api/game/{gameId}/stream/{userId} [get]
func StreamHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	userID := c.Param("userId")
	if gameID == "" || userID == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	ch := b.subscribe(gameID)
	defer b.unsubscribe(gameID, ch)

	var user model.User
	db.DB.First(&user, "id = ?", userID)
	var gameModel model.Game
	db.DB.First(&gameModel, "id = ?", gameID)
	db.DB.Model(&model.GamePlayer{}).Where("game_id = ? AND user_id = ?", gameID, userID).Update("is_connected", true)
	logAndSend(gameID, userID, "status", user.Name+": connected to the game")

	c.Stream(func(w io.Writer) bool {
		select {
		case ev := <-ch:
			data, _ := json.Marshal(ev)
			c.SSEvent("message", string(data))
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})

	db.DB.Model(&model.GamePlayer{}).Where("game_id = ? AND user_id = ?", gameID, userID).Update("is_connected", false)
	logAndSend(gameID, userID, "status", user.Name+": disconnected from the game")
	if gameModel.RequesterID == userID {
		db.DB.Model(&model.Game{}).Where("id = ?", gameID).Update("status", "abandoned")
		logAndSend(gameID, userID, "status", "Game was terminated because the creator disconnected")
		publishState(gameID)
	}
}

// ChatHandler records a chat message.
func ChatHandler(c *gin.Context) {
	var req struct {
		GameID  string `json:"gameId"`
		UserID  string `json:"userId"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.GameID == "" || req.UserID == "" || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	runes := []rune(msg)
	if len(runes) > 128 {
		msg = string(runes[:128])
	}
	msg = html.EscapeString(msg)
	var user model.User
	db.DB.First(&user, "id = ?", req.UserID)
	logAndSend(req.GameID, req.UserID, "chat", user.Name+": "+msg)
	c.Status(http.StatusOK)
}

// LogsHandler returns existing session logs for a game.
//
// @Summary      List session logs
// @Description  Retrieves chat and status logs for a game in reverse chronological order
// @Tags         events
// @Produce      json
// @Param        gameId  path  string  true  "Game ID"
// @Success      200  {array}  model.GameSessionLog
// @Router       /api/game/{gameId}/logs [get]
func LogsHandler(c *gin.Context) {
	gameID := c.Param("gameId")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing"})
		return
	}
	var logs []model.GameSessionLog
	db.DB.Where("game_id = ?", gameID).Order("created_at desc").Find(&logs)
	c.JSON(http.StatusOK, logs)
}

// LegacyAbandonHandler is the old abandon handler (deprecated, use AbandonGameHandler in game.go)
func LegacyAbandonHandler(c *gin.Context) {
	var req struct {
		GameID string `json:"gameId"`
		UserID string `json:"userId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.GameID == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	var gameModel model.Game
	if err := db.DB.First(&gameModel, "id = ?", req.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	if gameModel.RequesterID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only requester can abandon"})
		return
	}
	gameModel.Status = "abandoned"
	now := time.Now()
	gameModel.CompletedAt = &now
	db.DB.Save(&gameModel)
	var user model.User
	db.DB.First(&user, "id = ?", req.UserID)
	logAndSend(gameModel.ID, req.UserID, "status", "Game was abandoned by "+user.Name)
	publishState(gameModel.ID)
	c.Status(http.StatusOK)
}
