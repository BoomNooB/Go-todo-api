package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"API/auth"
	"API/todo"
)

var (
	buildCommit = "dev"
	buildTime   = time.Now().String()
)

func main() {
	//Liveness

	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Create("tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Printf("please provide your env file: %s", err)
	}
	fmt.Println(os.Getenv("SIGN"))
	fmt.Println(os.Getenv("DB_CONN"))

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("Failed to connect DB")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080", //Front-end path
	}

	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}

	r.Use(cors.New(config))

	//Readiness
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("limitz", limitedHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Build Commit": buildCommit,
			"Build Time":   buildTime,
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)

	// เพิ่ม get route เพื่อ handle หน้าบ้าน
	protected.GET("/todos", handler.List)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen :%s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("Shuttingdown gracefully, Press Ctrl+C again to FORCE")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}

}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
