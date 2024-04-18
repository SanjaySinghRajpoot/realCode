package main

import (
	"consumer/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CodeRunner struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type CodeRunnerRes struct {
	CodeResult    string `json:"codeResult"`
	CorrelationID string `json:"correlationID"`
}

const (
	PYTHON = "python"
	GOLANG = "golang"
)

func main() {
	// Set up configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafkaRealCode:19092", // Replace with your Kafka broker address
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	createTopics(context.Background(), []string{PYTHON, GOLANG}, config)

	producer, err := InitializeProducer()

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err.Error())
		return
	}

	// Create consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Subscribe to a topics
	topics := []string{PYTHON, GOLANG}
	consumer.SubscribeTopics(topics, nil)

	// Handle messages and shutdown signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	var codeObj CodeRunner
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		default:
			// time out of 100 millisecond
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				err = json.Unmarshal(e.Value, &codeObj)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, codeObj.Code)

				messageType := e.TopicPartition.Topic
				switch *messageType {
				case PYTHON:
					codeResult, err := utils.CompileCodePython(codeObj.Code)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(codeResult)
					sendResponse(PYTHON, codeResult, getCodeCorrelationID(e.Headers), producer)
				case GOLANG:
					codeResult, err := utils.CompileCodeGo(codeObj.Code)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(codeResult)
					sendResponse(GOLANG, codeResult, getCodeCorrelationID(e.Headers), producer)
				}

				// Send the response back to the producer

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				run = false

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func InitializeProducer() (*kafka.Producer, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafkaRealCode:19092",
		"client.id":         "my-group",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}

	return producer, nil
}

func getCodeCorrelationID(headers []kafka.Header) string {
	for _, header := range headers {
		if header.Key == "correlation-id" {
			return string(header.Value)
		}
	}
	return ""
}

func sendResponse(topic string, codeResult, correlationID string, producer *kafka.Producer) {
	// Convert response to bytes

	println("consusmer - producer setup")

	responseBytes, err := json.Marshal(map[string]string{
		"codeResult":    codeResult,
		"correlationID": correlationID,
	})
	if err != nil {
		log.Fatal(err)
	}

	resTopic := topic + "sec"

	// Produce the response message
	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &resTopic, Partition: kafka.PartitionAny},
		Value:          responseBytes,
		Headers: []kafka.Header{
			{Key: "correlation-id", Value: []byte(correlationID)},
		},
	}, deliveryChan)

	if err != nil {
		msg := fmt.Sprintf("Failed to produce message 1: %v\n", err)

		fmt.Println(msg)
		return

	} else {

		// Wait for delivery report
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			msg := fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error)
			fmt.Println(msg)
			return

		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}

		fls := producer.Flush(0) // Wait for up to 15 seconds for message delivery
		println(fls)
	}

}

func createTopics(ctx context.Context, topics []string, config *kafka.ConfigMap) {

	adminClient, err := kafka.NewAdminClient(config)
	if err != nil {
		panic(err)
	}
	defer adminClient.Close()

	topicSpecs := make([]kafka.TopicSpecification, len(topics))
	for i, topic := range topics {
		topicSpecs[i] = kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	results, err := adminClient.CreateTopics(ctx, topicSpecs, kafka.SetAdminOperationTimeout(5000))
	if err != nil {
		panic(err)
	}

	// Check if the topic creation was successful
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			panic(result.Error)
		}
	}
}