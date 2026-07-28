package usagelog

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type UsageLogRepository interface {
	SaveBatch(ctx context.Context, events []UsageLogEvent) error
}

type ClickHouseUsageLogRepository struct {
	conn clickhouse.Conn
}

func NewClickHouseUsageLogRepository(conn clickhouse.Conn) *ClickHouseUsageLogRepository {
	return &ClickHouseUsageLogRepository{conn: conn}
}

func (r *ClickHouseUsageLogRepository) SaveBatch(ctx context.Context, events []UsageLogEvent) error {
	if len(events) == 0 {
		return nil
	}

	batch, err := r.conn.PrepareBatch(ctx, `
		INSERT INTO usage_log
			(tenant, flag_name, enabled, called_at)
		VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare usage-log insert batch: %w", err)
	}

	for _, event := range events {
		if err := batch.Append(event.Tenant, event.FlagName, event.Enabled, event.CalledAt.UTC()); err != nil {
			return fmt.Errorf("append usage-log event: %w", err)
		}
	}
	if err := batch.Send(); err != nil {
		return fmt.Errorf("send usage-log insert batch: %w", err)
	}
	return nil
}
