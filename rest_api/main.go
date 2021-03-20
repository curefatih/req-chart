package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
)

var wg sync.WaitGroup
var KAFKA_TOPIC string = "apiRequests"
var kafkaProducer *kafka.Producer

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		panic(err)
	}

	kafkaProducer = p

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/", reqHandler).Methods(http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPut)

	// log dosyası ayarları
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	// log için çıktısı için dosya ve başına timestamp yazdırmaması için flag 0
	log.SetOutput(f)
	log.SetFlags(0)

	fmt.Println("Starting server...")

	http.ListenAndServe(":8080", r)

	wg.Wait()

	// Wait for message deliveries before shutting down
	kafkaProducer.Flush(15 * 1000)
}

func worker(val string) {
	kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &KAFKA_TOPIC, Partition: kafka.PartitionAny},
		Value:          []byte(val),
	}, nil)
}

func reqHandler(w http.ResponseWriter, r *http.Request) {

	// istek zamanı
	timeStart := time.Now().Unix()

	// random bekle
	n := rand.Intn(3)
	wait := time.Duration(n) * time.Second
	time.Sleep(wait)

	timeSpent := int64(wait / time.Millisecond)
	usedMethod := r.Method

	val := fmt.Sprintf("%s,%d,%d", usedMethod, timeSpent, timeStart)
	log.Printf(val)

	worker(val)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
