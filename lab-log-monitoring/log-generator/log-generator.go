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
	// Read environment variables
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	if broker == "" || topic == "" {
		fmt.Println(formatLog("ERROR", "Environment variables KAFKA_BROKER and KAFKA_TOPIC must be set"))
		os.Exit(1)
	}

	// Set up Kafka writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	rand.Seed(time.Now().UnixNano())

	for {
		level := getRandomLogLevel()
		message := fmt.Sprintf("Log entry at %s", time.Now().Format(time.RFC3339))
		logLine := formatLog(level, message)

		// Always print to stdout
		fmt.Println(logLine)

		// Try sending to Kafka, do not crash on failure
		err := writer.WriteMessages(nil, kafka.Message{
			Value: []byte(logLine),
		})
		if err != nil {
			// Log Kafka error in the same format
			errorLog := formatLog("ERROR", fmt.Sprintf("Kafka write failed: %v", err))
			fmt.Println(errorLog)
		}

		time.Sleep(1 * time.Second)
	}
}
