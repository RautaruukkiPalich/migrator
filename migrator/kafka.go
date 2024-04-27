package migrator

import (
	"context"
	"fmt"
	"sync/atomic"

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
	data, err := m.getMsgsFromRows(table, sqlrows)
	if err != nil {
		return err
	}

	err = m.broker.WriteMessages(
		context.Background(),
		data...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *migrator) getMsgsFromRows(table string, rows *sqlx.Rows) ([]kafka.Message, error) {
	headers, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var cursor int64 = 0
	data := make([]kafka.Message, batchSize)

	for rows.Next() {
		row, _ := rows.SliceScan()

		// закинемданные в мапу, которую и отправим, чтобы не было проблем с именами столбцов 
		dict := make(map[string]any)
		for idx, col := range headers {
			dict[col] = row[idx]
		} 

		msg := kafka.Message{
			// возможно, стоит название таблицы записать в key, а не в header
			Headers: []kafka.Header{
				{
					Key:   "table",
					Value: []byte(table),
				},
			},
			// возможно, стоит данные сериализовать и отправить в Value
			WriterData: dict,
		}
		data[cursor] = msg
		atomic.AddInt64(&cursor, 1)
	}

	return data[:cursor], nil
}
