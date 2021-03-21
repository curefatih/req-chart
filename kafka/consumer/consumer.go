package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	Type      string
	Spent     int64
	Timestamp int64
}

func main() {

	dsn := "host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connected to database", db)
	}

	db.AutoMigrate(&Request{})

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka",
		"group.id":          "consumer1",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		fmt.Println(err)

		// retry connection
		for err != nil {

			time.Sleep(2 * time.Second)

			c, err = kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": "kafka",
				"group.id":          "consumer1",
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				fmt.Println(err)
			}
		}
	}

	c.SubscribeTopics([]string{"apiRequests"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s", msg.TopicPartition, string(msg.Value))

			val := strings.Split(string(msg.Value), ",")
			fmt.Println("vals:", val)
			if len(val) != 3 {
				fmt.Println("Passing...")
				continue
			}

			spent, err := strconv.Atoi(val[1])
			if err != nil {
				fmt.Println("Cant parse time spent.", err)
				continue
			}

			timestamp, err := strconv.Atoi(val[2])
			if err != nil {
				fmt.Println("Cant parse timestamp", err)
				continue
			}

			db.Create(&Request{
				Type:      val[0],
				Spent:     int64(spent),
				Timestamp: int64(timestamp),
			})

		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
