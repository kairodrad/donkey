package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/model"
)

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
