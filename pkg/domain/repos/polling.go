package repos

import (
	"polling_websocket/pkg/domain/models"
	"time"
)

type PollingService interface {
	GetContentActionByID(actionID *string, userID *string) (data *string, err error)
}

type PollingHTTPRepository interface {
	GetActionByID(actionID *string, userID *string, commandType string, limitCount uint64) (data *models.ResponsePollingActionID, err error)
}

type PollingRedisRepoInterface interface {
	ValidateActionGlobalUUID(field *string) (bool, error)
	SetNX(hashKey, actionID string, expiration time.Duration) (bool, error)
}

type PollingBrokerRepository interface {
	SendAction(newAction *models.RequestGoogleAction) bool
}

type CredentialBrokerRepository interface {
	UpdateCredential(exchangeCredential *models.RequestExchangeCredential) bool
}
