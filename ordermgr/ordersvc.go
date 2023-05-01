package ordermgr

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserOrder(userId string) (Order, error) {
	fmt.Printf("received request to get order for user %s, \n", userId)
	var order Order
	filter := bson.D{{"userid", userId}}
	err := orderCollection.FindOne(context.Background(), filter).Decode(&order)
	if err != nil {
		fmt.Printf("failed to read order for user: %s, error: %v\n", userId, err)
		return order, err
	}
	return order, nil
}

type Order struct {
	UserId string `json:"userid,omitempty" bson:"userid,omitempty"`
	ID     string `json:"id"`
	Price  string `json:"price"`
	Date   string `json:"date"`
}

func CreateUserOrder(userId string) error {
	fmt.Printf("received request to creare new order for user %s, \n", userId)
	id, price := genOrderIdAndPrice()
	order := Order{
		UserId: userId,
		ID:     id,
		Price:  price,
		Date:   time.Now().Format(time.RFC3339),
	}
	result, err := orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		fmt.Printf("failed to create order for user: %s, error: %v\n", userId, err)
		return err
	}
	fmt.Printf("user %s, order created with InsertedID: %v\n", userId, result.InsertedID)
	return nil
}

func genOrderIdAndPrice() (string, string) {
	// Create a big.Int with the maximum value for the desired range
	max := big.NewInt(10000)
	randInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Println("Error generating random order number:", err)
		return "100", "1000.00"
	}
	id := randInt.String()
	price := randInt.Add(randInt, big.NewInt(102))
	priceFloat := big.NewFloat(float64(price.Int64())).String()
	return id, priceFloat
}
