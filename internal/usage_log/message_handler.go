package usagelog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	usageLogTopic = "feature-flags-call-log"
	batchSize     = 8192
	batchWait     = 60 * time.Second
	readTimeout   = 250 * time.Millisecond
)

type BatchProcessor func(context.Context, []*kafka.Message) error

func RegisterUsageLogHandler(ctx context.Context, process BatchProcessor) {
	go func() {
		if err := RunUsageLogHandler(ctx, process); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("Usage log consumer stopped", "error", err)
		}
	}()
}

// RunUsageLogHandler consumes usage events and gives them to process in batches.
// It flushes as soon as either 8,192 messages are buffered or 60 seconds have
// elapsed since the first message entered the current batch.
func RunUsageLogHandler(ctx context.Context, process BatchProcessor) error {
	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		bootstrapServers = "localhost:9092"
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           "usage-log-consumer",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return fmt.Errorf("create Kafka consumer: %w", err)
	}
	defer consumer.Close()

	if err := consumer.SubscribeTopics([]string{usageLogTopic}, nil); err != nil {
		return fmt.Errorf("subscribe to %q: %w", usageLogTopic, err)
	}

	batch := make([]*kafka.Message, 0, batchSize)
	// deadline is zero while there is no active batch timer.
	var deadline time.Time

	flush := func(processingCtx context.Context) error {
		if len(batch) == 0 {
			return nil
		}

		if err := process(processingCtx, batch); err != nil {
			return fmt.Errorf("process usage-log batch: %w", err)
		}

		offsets := make(map[kafka.TopicPartition]kafka.Offset)
		for _, msg := range batch {
			partition := kafka.TopicPartition{
				Topic:     msg.TopicPartition.Topic,
				Partition: msg.TopicPartition.Partition,
			}
			nextOffset := msg.TopicPartition.Offset + 1
			if offset, ok := offsets[partition]; !ok || nextOffset > offset {
				offsets[partition] = nextOffset
			}
		}

		commitOffsets := make([]kafka.TopicPartition, 0, len(offsets))
		for partition, offset := range offsets {
			partition.Offset = offset
			commitOffsets = append(commitOffsets, partition)
		}
		if _, err := consumer.CommitOffsets(commitOffsets); err != nil {
			return fmt.Errorf("commit usage-log offsets: %w", err)
		}

		// Reuse the allocated storage, but remove references to processed messages.
		clear(batch)
		batch = batch[:0]
		deadline = time.Time{}
		return nil
	}

	for {
		if err := ctx.Err(); err != nil {
			// A partial batch is still useful on graceful shutdown. It is committed
			// only if processing succeeds; otherwise Kafka will redeliver it.
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if flushErr := flush(shutdownCtx); flushErr != nil {
				return flushErr
			}
			return err
		}

		// Check before reading again so a busy topic cannot starve the time flush.
		if !deadline.IsZero() && !time.Now().Before(deadline) {
			if err := flush(ctx); err != nil {
				return err
			}
			continue
		}

		msg, err := consumer.ReadMessage(readTimeout)
		if err != nil {
			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && kafkaErr.IsTimeout() {
				continue // The short read timeout lets us notice context and timer expiry.
			}
			return fmt.Errorf("read Kafka message: %w", err)
		}

		if len(batch) == 0 {
			// The 60-second window starts with the first message, not on a global tick.
			deadline = time.Now().Add(batchWait)
		}
		batch = append(batch, msg)

		if len(batch) >= batchSize {
			if err := flush(ctx); err != nil {
				return err
			}
		}
	}
}
