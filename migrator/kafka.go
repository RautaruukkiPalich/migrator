package migrator

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	timeout = time.Millisecond * 1
	topic = "migrator"
)

type broker struct {
	producer *kafka.Writer
	consumer *kafka.Reader
	ticker *time.Ticker
	done chan struct{}
}

func newBroker(cfg *KafkaConfig) (*broker, error) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)), //"localhost:29092")
		Balancer: &kafka.LeastBytes{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Topic:     topic,
		MaxBytes:  10e6, // 10MB
	})

	ticker := time.NewTicker(timeout)
	done := make(chan struct{})

	return &broker{
		producer: w,
		consumer: r,
		ticker: ticker,
		done: done,
	}, nil
}

func (b *broker) Run() {
	
	ctx := context.Background()
	
	go func() {
		for {
			select {
			case <- b.done:
				return
			case <- b.ticker.C:
				m, err := b.consumer.ReadMessage(ctx)
				if err != nil {
					fmt.Printf("error reading message %v", err)
				}
				fmt.Println(m)
			}
		}
	}()
}

func (b *broker) Stop(){
	close(b.done)
	defer b.ticker.Stop()
	defer b.producer.Close()
	defer b.consumer.Close()
}

func (b *broker) SendMessages() {
	data := "message"
	err := b.producer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: topic,
			Value: []byte(data),
		})
	if err != nil {
		panic(err) 
	}
}

func (b *broker) GetMessages() {
	m, err := b.consumer.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}