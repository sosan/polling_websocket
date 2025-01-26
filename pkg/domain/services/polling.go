package services

import (
	"fmt"
	"polling_websocket/pkg/domain/models"
	"polling_websocket/pkg/domain/repos"
)

type PollingServiceImpl struct {
	redisRepo          repos.PollingRedisRepoInterface
	brokerPollingRepo  repos.PollingBrokerRepository
	brokerCredsRepo    repos.CredentialBrokerRepository
	httpRepo           repos.PollingHTTPRepository
	credentialHTTPRepo repos.CredentialHTTPRepository
}

func NewPollingService(
	repoRedis repos.PollingRedisRepoInterface,
	actionBroker repos.PollingBrokerRepository,
	repoHTTP repos.PollingHTTPRepository,
	credentialRepo repos.CredentialHTTPRepository,
	credentialBroker repos.CredentialBrokerRepository,
) repos.PollingService {
	return &PollingServiceImpl{
		redisRepo:          repoRedis,
		brokerPollingRepo:  actionBroker,
		brokerCredsRepo:    credentialBroker,
		httpRepo:           repoHTTP,
		credentialHTTPRepo: credentialRepo,
	}
}

func (a *PollingServiceImpl) GetContentActionByID(actionID *string, userID *string) (data *string, err error) {
	if actionID == nil || userID == nil {
		return nil, fmt.Errorf("cannot by empty")
	}

	data, err = a.getAllContentFromAction(actionID, userID, models.UpdateCommand)
	return data, err
}
