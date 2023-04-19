package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Order struct {
	ID    string `json:"id"`
	Price string `json:"price"`
	Date  string `json:"date"`
}

func main() {
	r := gin.Default()
	r.GET("/order", GetUserHandler)
	err := r.Run(":9092")
	if err != nil {
		log.Fatalf("impossible to start server: %s", err)
	}
}

func GetUserHandler(c *gin.Context) {
	id := c.Query("id")
	teacher, err := GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve order"})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func GetOrderByID(id string) (Order, error) {
	// TODO : lookup in db
	return Order{
		ID:    id,
		Price: "1000.00",
		Date:  time.Now().Format(time.RFC3339),
	}, nil
}
