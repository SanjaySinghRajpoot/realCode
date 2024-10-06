package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SanjaySinghRajpoot/realCode/backend/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var KafkaProducer *kafka.Producer

func Producer(Topic string, post models.CodeRunner, producer *kafka.Producer) (string, error) {

	// Generate a unique correlation ID
	correlationID := uuid.NewV4().String()

	// Convert struct to bytes
	postBytes, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &Topic, Partition: kafka.PartitionAny},
		Value:          postBytes,
		Headers: []kafka.Header{
			{Key: "correlation-id", Value: []byte(correlationID)},
		},
	}, deliveryChan)

	if err != nil {
		msg := fmt.Sprintf("Failed to produce message 1: %v\n", err)

		return msg, err

	} else {

		// Wait for delivery report
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			msg := fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error)

			return msg, m.TopicPartition.Error

		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}

		fls := producer.Flush(0) // Wait for up to 15 seconds for message delivery
		println(fls)
	}

	// setup consumer here
	// Set up configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafkaRealCode:19092", // Replace with your Kafka broker address
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	// Create consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Subscribe to a topics
	topics := []string{"pythonsec", "golangsec"}
	consumer.SubscribeTopics(topics, nil)

	createTopics(context.Background(), topics, config)

	// Handle messages and shutdown signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	var codeObj models.CodeRunnerRes
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

				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, codeObj.CorrelationID)

				if codeObj.CorrelationID == correlationID {
					return codeObj.CodeResult, nil
				}

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				run = false

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	return "test", nil
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

	fmt.Println("producer created")

	return producer, nil
}

func getCorrelationID(headers []kafka.Header) string {
	for _, header := range headers {
		if header.Key == "correlation-id" {
			return string(header.Value)
		}
	}
	return ""
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
