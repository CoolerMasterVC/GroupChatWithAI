package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"transport/internal/models"
	"transport/internal/storage"

	"github.com/IBM/sarama"
)

var (
	brokers    []string
	topic      = "segments"
	consumerCx context.Context
	cancelFunc context.CancelFunc
)

func init() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	brokers = strings.Split(broker, ",")
}

// Producer
func ProduceSegment(segment models.TransferRequest) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer producer.Close()

	data, err := json.Marshal(segment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	_, _, err = producer.SendMessage(msg)
	return err
}

// Consumer
func StartConsumer(store *storage.Storage) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}

	consumerCx, cancelFunc = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				var seg models.TransferRequest
				if err := json.Unmarshal(msg.Value, &seg); err != nil {
					log.Printf("Unmarshal error: %v", err)
					continue
				}
				// Валидация номера сегмента
				if seg.SegmentNumber < 1 || seg.SegmentNumber > seg.TotalSegments {
					log.Printf("[KAFKA] invalid segment number %d (total %d) from %s, skipping", seg.SegmentNumber, seg.TotalSegments, seg.SendTime)
					continue
				}
				store.AddOrUpdate(seg.SendTime, seg.SegmentNumber, seg.TotalSegments, seg.Username, seg.Payload)
				log.Printf("[KAFKA] consumed segment %d/%d for %s", seg.SegmentNumber, seg.TotalSegments, seg.SendTime)

			case err := <-partitionConsumer.Errors():
				log.Printf("Consumer error: %v", err)

			case <-consumerCx.Done():
				partitionConsumer.Close()
				consumer.Close()
				return
			}
		}
	}()
}

func StopConsumer() {
	if cancelFunc != nil {
		cancelFunc()
	}
}
