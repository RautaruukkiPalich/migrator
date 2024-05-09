package migrator

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/rautaruukkipalich/migrator/config"
	"github.com/segmentio/kafka-go"
)

func newBroker(cfg *config.KafkaConfig) *kafka.Writer {

	// hint for fix: panic: [3] Unknown Topic Or Partition: the request is for a topic or partition that does not exist on this broker
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Topic,
		0,
	)
	if err != nil {
		panic(err)
	}
	// close the connection because we won't be using it
	conn.Close()

	w := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)), //"localhost:29092"
		Topic:    Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return w
}

func (m *migrator) SendMessages(table string, rows *sqlx.Rows) error {
	defer rows.Close()

	msgs, err := m.getMsgsFromRows(table, rows)
	if err != nil {
		return err
	}

	err = m.broker.WriteMessages(
		context.Background(),
		msgs...,
	)
	if err != nil {
		return ErrFailedToSendKafkaMessages
	}

	return nil
}

func (m *migrator) getMsgsFromRows(table string, rows *sqlx.Rows) ([]kafka.Message, error) {

	var cursor int32

	msgs := make([]kafka.Message, m.batchSize)

	for rows.Next() {
		rowMap := make(map[string]any)

		if err := rows.MapScan(rowMap); err != nil {
			return nil, err
		}

		value, err := jsoniter.Marshal(rowMap)
		if err != nil {
			return nil, ErrFailedToMarshal
		}

		key := []byte(fmt.Sprintf("%s_%d", table, rowMap["id"]))

		msgs[cursor] = createMsg(table, key, value)
		atomic.AddInt32(&cursor, 1)
	}

	return msgs[:cursor], nil
}

func createMsg(table string, key, value []byte) kafka.Message {
	return kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "table",
				Value: []byte(table),
			},
		},
		Key:   key,
		Value: value,
	}
}
