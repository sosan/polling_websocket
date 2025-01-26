package brokerclient

import (
	"encoding/json"
	// "fmt"
	"log"
	"polling_websocket/pkg/common"
	"polling_websocket/pkg/domain/models"
	"time"
)

type CredentialCommand struct {
	Type       string                            `json:"type,omitempty"`
	Credential *models.RequestExchangeCredential `json:"credential"`
	Timestamp  time.Time                         `json:"timestamp,omitempty"`
}

type CredentialKafkaRepository struct {
	client KafkaClient
}

func NewCredentialKafkaRepository(client KafkaClient) *CredentialKafkaRepository {
	return &CredentialKafkaRepository{
		client: client,
	}
}

func (c *CredentialKafkaRepository) UpdateCredential(payload *models.RequestExchangeCredential) (sended bool) {
	if payload == nil {
		log.Printf("ERROR | Payload is nil")
		return false
	}
	command := CredentialCommand{
		Credential: payload,
		Type: models.UpdateCommand,
	}
	sended = c.PublishCommand(command, payload.ID)
	return sended
}

func (c *CredentialKafkaRepository) PublishCommand(credentialCommand CredentialCommand, key string) bool {
	command, err := json.Marshal(credentialCommand)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = c.client.Produce(models.CredentialTopicName, []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.Sleep(waitTime)
	}

	return false
}
