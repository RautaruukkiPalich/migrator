package migrator

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

const (
	topic = "migrator"
)

func newBroker(cfg *KafkaConfig) *kafka.Writer {
	// hint for fix: panic: [3] Unknown Topic Or Partition: the request is for a topic or partition that does not exist on this broker
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, 0)
	if err != nil {
		panic(err)
	}
	// close the connection because we won't be using it
	conn.Close()

	w := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)), //"localhost:29092")
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return w
}

func (m *migrator) SendMessages(table string, rows *sqlx.Rows) error {
	sqlrows := rows
	data := m.getMsgsFromRows(table, sqlrows)

	err := m.broker.WriteMessages(
		context.Background(),
		data...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *migrator) getMsgsFromRows(table string, rows *sqlx.Rows) []kafka.Message {
	var data []kafka.Message
	for rows.Next() {
		row, _ := rows.SliceScan()
		msg := kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   "table",
					Value: []byte(table),
				},
			},
			WriterData: row,
		}
		data = append(data, msg)
	}

	return data
}
