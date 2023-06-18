package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"API/auth"
	"API/todo"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect DB")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/tokenz", auth.AccessToken)

	handler := todo.NewTodoHandler(db)
	r.POST("/todos", handler.NewTask)

	r.Run()
}
