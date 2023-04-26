package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID    string `json:"id"`
	Price string `json:"price"`
	Date  string `json:"date"`
}

func main() {
	r := gin.Default()

	serverPort := GetEnvParam("SERVICE_PORT", "8081")

	r.GET("/order", GetUserHandler)
	err := r.Run(":" + serverPort)
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

// GetEnvParam : return string environmental param if exists, otherwise return default
func GetEnvParam(param string, dflt string) string {
	if v, exists := os.LookupEnv(param); exists {
		return v
	}
	return dflt
}
