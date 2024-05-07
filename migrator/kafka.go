package migrator

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

const (
	topic = "migrator"
)

func newBroker(cfg *KafkaConfig) *kafka.Writer {

	// hint for fix: panic: [3] Unknown Topic Or Partition: the request is for a topic or partition that does not exist on this broker
	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), topic, 0)
	if err != nil {
		panic(err)
	}
	// close the connection because we won't be using it
	conn.Close()

	w := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)), //"localhost:29092"
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return w
}

func (m *migrator) SendMessages(table string, rows *sqlx.Rows) error {
	defer rows.Close()

	data, err := m.getMsgsFromRows(table, rows)
	if err != nil {
		return err
	}

	err = m.broker.WriteMessages(
		context.Background(),
		data...,
	)
	if err != nil {
		return ErrFailedToSendKafkaMessages
	}

	return nil
}

func (m *migrator) getMsgsFromRows(table string, rows *sqlx.Rows) ([]kafka.Message, error) {

	var (
		cursor int64 = 0
		key    []byte
		value  []byte
		err    error
	)

	data := make([]kafka.Message, batchSize)

	for rows.Next() {
		switch table {
		case "donor":
			type DBUser struct {
				ID        int64     `json:"id"`
				Username  string    `json:"username"`
				CreatedAt time.Time `json:"created_at"`
			}

			user := DBUser{}
			rows.Scan(&user.ID, &user.Username, &user.CreatedAt)

			value, err = jsoniter.Marshal(user)
			if err != nil {
				return nil, ErrFailedToMarshal
			}

			key = []byte(fmt.Sprintf("%s_%d", table, user.ID))

			//или не зная структуры, использовать рефлект?
		default:
			//что то тут для других таблиц
		}

		msg := kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   "table",
					Value: []byte(table),
				},
			},
			Key:   key,
			Value: value,
		}
		data[cursor] = msg
		atomic.AddInt64(&cursor, 1)

	}

	return data[:cursor], nil
}
