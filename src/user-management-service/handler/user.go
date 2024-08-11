package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	model "filmsharing/user-management-service/model"
)

var UserMemDB = []model.User{
	{ID: "1", Name: "John Doe", Email: "john.doe@example.com", Password: "password", Username: "john.doe"},
	{ID: "2", Name: "Jane Doe", Email: "jane.doe@example.com", Password: "password", Username: "jane.doe"},
	{ID: "3", Name: "John Smith", Email: "john.smith@example.com", Password: "password", Username: "john.smith"},
	{ID: "4", Name: "Jane Smith", Email: "jane.smith@example.com", Password: "password", Username: "jane.smith"},
	{ID: "5", Name: "John Doe", Email: "john.doe@example.com", Password: "password", Username: "john.doe"},
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, UserMemDB)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, user := range UserMemDB {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	UserMemDB = append(UserMemDB, user)
	c.JSON(http.StatusAccepted, user)
}
