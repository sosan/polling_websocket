package brokerclient

import (
	"encoding/json"
	"log"
	"polling_websocket/pkg/common"
	"polling_websocket/pkg/domain/models"
	"time"
)

type PollingKafkaRepository struct {
	client KafkaClient
}

const (
	CommandTypeCreate = "create"
	CommandTypeUpdate = "update"
	CommandTypeDelete = "delete"
	TopicName         = "actions.command"
)

type PollingCommand struct {
	Polling   *models.RequestGoogleAction `json:"actions"`
	Type      string                      `json:"type,omitempty"`
	Timestamp time.Time                   `json:"timestamp,omitempty"`
}

func NewPollingKafkaRepository(client KafkaClient) *PollingKafkaRepository {
	return &PollingKafkaRepository{
		client: client,
	}
}

func (a *PollingKafkaRepository) SendAction(newAction *models.RequestGoogleAction) (sended bool) {
	command := PollingCommand{
		Polling:   newAction,
		Type:      CommandTypeUpdate,
		Timestamp: time.Now(),
	}
	sended = a.PublishCommand(command, newAction.ActionID)
	return sended
}

func (a *PollingKafkaRepository) PublishCommand(payload PollingCommand, key string) bool {
	command, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = a.client.Produce(TopicName, []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.After(waitTime)
	}

	return false
}
