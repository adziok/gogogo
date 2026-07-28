package usagelog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type UsageLogEvent struct {
	Tenant   string
	FlagName string
	Enabled  bool
	CalledAt time.Time
}

type usageLogMessage struct {
	Tenant    string    `json:"tenant"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func NewClickHouseProcessor(repository UsageLogRepository) BatchProcessor {
	return func(ctx context.Context, messages []*kafka.Message) error {
		events := make([]UsageLogEvent, 0, len(messages))
		for _, message := range messages {
			var payload usageLogMessage
			if err := json.Unmarshal(message.Value, &payload); err != nil {
				return fmt.Errorf("decode Kafka message at %s: %w", message.TopicPartition, err)
			}
			events = append(events, UsageLogEvent{
				Tenant: payload.Tenant, FlagName: payload.Name, Enabled: payload.Enabled, CalledAt: payload.CreatedAt,
			})
		}
		if err := repository.SaveBatch(ctx, events); err != nil {
			return fmt.Errorf("save usage-log batch: %w", err)
		}
		return nil
	}
}
