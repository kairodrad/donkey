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

// RenameHandler updates a user's display name. This endpoint is used internally and is
// intentionally undocumented in Swagger.
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

// GetUserHandler returns a user by ID.
//
// @Summary      Get user
// @Tags         user
// @Produce      json
// @Param        id  path  string  true  "user id"
// @Success      200  {object}  model.User
// @Failure      404  {object}  map[string]string
// @Router       /api/user/{id} [get]
func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.DB.Preload("Games").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ListUsersHandler returns all registered users.
//
// @Summary      List users
// @Tags         user
// @Produce      json
// @Success      200  {array}  model.User
// @Router       /api/users [get]
func ListUsersHandler(c *gin.Context) {
	var users []model.User
	db.DB.Preload("Games").Find(&users)
	c.JSON(http.StatusOK, users)
}
