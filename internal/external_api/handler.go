package externalapi

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"start/internal/auth"
	"start/internal/common"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-chi/chi/v5"
)

type ExternalApiHandler struct {
	repository  FeatureFlagExternalReposiotory
	emitMessage func(msg []byte) error
	close       func()
}

func CreateExternalApiHandler(repo FeatureFlagExternalReposiotory) (*ExternalApiHandler, error) {
	emitMessage, closeProducer, err := createKafkaProducer()
	if err != nil {
		return nil, fmt.Errorf("create Kafka producer: %w", err)
	}

	return &ExternalApiHandler{
		repository:  repo,
		emitMessage: emitMessage,
		close:       closeProducer,
	}, nil
}

// Close flushes and closes the Kafka producer. Call it once during server shutdown.
func (h *ExternalApiHandler) Close() {
	if h.close != nil {
		h.close()
		h.close = nil
	}
}

func (h *ExternalApiHandler) GetByTenantAndName(w http.ResponseWriter, r *http.Request) {
	flagName := chi.URLParam(r, "id")
	userDetails := auth.GetUserDetails(r.Context())

	data, err := h.repository.GetByTenantAndName(r.Context(), userDetails.OrgID, flagName)
	if err != nil {
		common.RenderErr(w, r, http.StatusInternalServerError, err, "Internal server error")
		return
	}
	go func() {
		msg := messageToKafka{
			Tenant:    userDetails.OrgID,
			Name:      flagName,
			Enabled:   data.Enabled,
			CreatedAt: time.Now(),
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			slog.Error("Failed to marshal Kafka message", "error", err)
			return
		}
		if emitError := h.emitMessage(jsonData); emitError != nil {
			slog.Error("Failed to publish Kafka message", "error", emitError)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}

}

type messageToKafka struct {
	Tenant    string    `json:"tenant"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func createKafkaProducer() (func(msg []byte) error, func(), error) {
	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		bootstrapServers = "localhost:9092"
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					slog.Error("Delivery failed: \n", ev.TopicPartition)
				} else {
					slog.Debug("Delivered message to \n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "feature-flags-call-log"

	return func(msg []byte) error {
			return p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          msg,
			}, nil)
		}, func() {
			p.Flush(5_000)
			p.Close()
		}, nil
}
