package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"API/auth"
	"API/todo"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Println("please provide your env file: %s", err)
	}
	fmt.Println(os.Getenv("SIGN"))

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

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)

	r.Run()
}
