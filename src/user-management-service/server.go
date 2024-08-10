package main

import (
	"github.com/gin-gonic/gin"

	handler "filmsharing/user-management-service/handler"
)

func main() {
	router := gin.Default()

	router.GET("/users", handler.GetUsers)
	router.GET("/users/:id", handler.GetUserByID)
	router.POST("/users", handler.CreateUser)

	router.Run(":8080")
}
