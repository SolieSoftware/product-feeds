package main

import (
	"log"

	"github.com/SolieSoftware/product-feeds/db"
	"github.com/SolieSoftware/product-feeds/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Failed to load .env: ", err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Private-Network", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.GET("/api/products", handlers.GetProducts)
	r.Run(":8080")
}
