package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

// RenameRequest is the request body for renaming a user.
type RenameRequest struct {
	GameID string `json:"gameId"`
	Name   string `json:"name"`
}

// RenameHandler updates a user's display name.
//
// @Summary      Rename user
// @Description  Update a user's name and notify the game if provided
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path  string         true  "user id"
// @Param        data  body  RenameRequest  true  "rename request"
// @Success      200  {object}  model.User
// @Router       /api/user/{id}/rename [post]
func RenameHandler(c *gin.Context) {
	userID := c.Param("id")
	var req RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil || userID == "" || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	var user model.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	old := user.Name
	user.Name = req.Name
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if req.GameID != "" {
		logAndSend(req.GameID, userID, "status", old+": Player renamed to "+user.Name)
		publishState(req.GameID)
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name})
}
