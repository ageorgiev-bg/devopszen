package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var logLevels = []string{"INFO", "WARNING", "ERROR"}

func getRandomLogLevel() string {
	return logLevels[rand.Intn(len(logLevels))]
}

func formatLog(level, message string) string {
	timestamp := time.Now().Format(time.RFC3339)
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	var writer *kafka.Writer
	kafkaEnabled := false

	if broker == "" || topic == "" {
		fmt.Println(formatLog("ERROR", "KAFKA_BROKER or KAFKA_TOPIC not set. Kafka disabled."))
	} else {
		writer = &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
		defer writer.Close()
		kafkaEnabled = true
	}

	rand.Seed(time.Now().UnixNano())

	for {
		level := getRandomLogLevel()
		message := fmt.Sprintf("Log entry at %s", time.Now().Format(time.RFC3339))
		logLine := formatLog(level, message)

		// Always print to stdout
		fmt.Println(logLine)

		// Send to Kafka only if enabled
		if kafkaEnabled {
			err := writer.WriteMessages(nil, kafka.Message{Value: []byte(logLine)})
			if err != nil {
				errMsg := formatLog("ERROR", fmt.Sprintf("Kafka write failed: %v", err))
				fmt.Println(errMsg)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
