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

func (m *migrator) SendToBroker(table string, rows *sqlx.Rows) error {

	defer rows.Close()

	var cursor int32
	var batch int32

	msgs := make([]kafka.Message, m.batchSize)

	for rows.Next() {
		rowMap := make(map[string]any)

		if err := rows.MapScan(rowMap); err != nil {
			return ErrParseRow
		}

		value, err := jsoniter.Marshal(rowMap)
		if err != nil {
			return ErrFailedToMarshal
		}

		key := []byte(fmt.Sprintf("%s_%d", table, batch * int32(m.batchSize) + cursor))

		msgs[cursor] = createMsg(table, key, value)
		atomic.AddInt32(&cursor, 1)

		if cursor == m.batchSize {
			atomic.AddInt32(&batch, 1)
			if err = m.SendMessages(msgs); err != nil {
				return err
			}

			msgs = make([]kafka.Message, m.batchSize)
			atomic.StoreInt32(&cursor, 0)
		}
	}

	return m.SendMessages(msgs[:cursor])
}

func (m *migrator) SendMessages(msgs []kafka.Message) error {
	err := m.broker.WriteMessages(
		context.Background(),
		msgs...,
	)
	if err != nil {
		return ErrFailedToSendKafkaMessages
	}

	return nil
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
