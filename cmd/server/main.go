package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/charly3pins/transferrer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := flag.String("dsn", os.Getenv("POSTGRES_DSN"), "db connection string")
	flag.Parse()

	db, err := transferrer.NewDB(*dsn)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	store := transferrer.NewStore(db)
	handler := transferrer.Handler{Store: store}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET, POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/token", transferrer.GenerateJWT)
	router.GET("/balance", transferrer.ValidateToken, handler.Balance)
	router.POST("/transfer", transferrer.ValidateToken, handler.Transfer)

	router.Run(":8080")
}
