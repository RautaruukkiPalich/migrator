package kafka

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/rautaruukkipalich/migrator/config"
	"github.com/rautaruukkipalich/migrator/migrator/broker"
	"github.com/rautaruukkipalich/migrator/migrator/database"
	"github.com/segmentio/kafka-go"
	jsoniter "github.com/json-iterator/go"
)

type Broker struct {
	writer *kafka.Writer
	batchSize int32
}

func NewBroker(cfg *config.Config) (*Broker, error) {

	// hint for fix: panic: [3] Unknown Topic Or Partition: the request is for a topic or partition that does not exist on this broker
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		cfg.Kafka.URI,
		cfg.Kafka.Topic,
		0,
	)
	if err != nil {
		return nil, err
	}
	// close the connection because we won't be using it
	conn.Close()

	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.URI), //"localhost:29092"
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Broker{
		writer: w,
		batchSize: cfg.BatchSize,
		}, nil
}

func (b *Broker) Close() {
	b.writer.Close()
}

func (b *Broker) SendMessages(ctx context.Context, rowsch chan database.Row, tablename string) error {
	var cursor int32
	var batch int32
	var msgs []kafka.Message

	resetCounter := func(){
		msgs = make([]kafka.Message, b.batchSize)
		atomic.StoreInt32(&cursor, 0)
	} 

	resetCounter()

	for row := range rowsch{

		if row.Err != nil {
			return row.Err
		}

		value, err := jsoniter.Marshal(row.Row)
		if err != nil {
			return broker.ErrFailedToMarshal
		}

		key := []byte(fmt.Sprintf("%s_%d", tablename, batch * b.batchSize + cursor))

		msgs[cursor] = createMsg(tablename, key, value)
		atomic.AddInt32(&cursor, 1)

		if cursor == b.batchSize {
			if err = b.Send(ctx, msgs); err != nil {
				return err
			}

			atomic.AddInt32(&batch, 1)
			resetCounter()
		}
	}

	return b.Send(ctx, msgs[:cursor])
}

func (b *Broker) Send(ctx context.Context, msgs []kafka.Message) error {
	err := b.writer.WriteMessages(
		ctx,
		msgs...,
	)
	if err != nil {
		return broker.ErrFailedToSendKafkaMessages
	}

	return nil
}

func createMsg(tablename string, key, value []byte) kafka.Message {
	return kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "table",
				Value: []byte(tablename),
			},
		},
		Key:   key,
		Value: value,
	}
}