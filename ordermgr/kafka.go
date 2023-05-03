package ordermgr

// import (
// 	"fmt"

// 	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
// )

// func Consumer() {
// 	fmt.Println("system service is ready to read from Kafka...")
// 	host := GetEnvParam("KAFKA_HOST", "localhost")
// 	port := GetEnvParam("KAFKA_PORT", "9092")
// 	topic := GetEnvParam("KAFKA_TOPIC", "order")

// 	fmt.Printf("Kafka host:%s , ,port:%s \n", host, port)

// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:9092",
// 		"group.id":          "orderGroup",
// 		"auto.offset.reset": "earliest",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	c.SubscribeTopics([]string{topic}, nil)

// 	// A signal handler or similar could be used to set this to false to break the loop.
// 	run := true

// 	for run {
// 		msg, err := c.ReadMessage(-1)
// 		if err == nil {
// 			fmt.Printf("Received message from %s: %s\n", msg.TopicPartition, string(msg.Value))

// 			CreateUserOrder(string(msg.Value))

// 		} else if !err.(kafka.Error).IsTimeout() {
// 			// The client will automatically try to recover from all errors.
// 			// Timeout is not considered an error because it is raised by
// 			// ReadMessage in absence of messages.
// 			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 		}
// 	}

// 	c.Close()
// }
