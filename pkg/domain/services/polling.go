package services

import (
	"fmt"
	"log"
	"polling_websocket/pkg/domain/repos"
)

type PollingServiceImpl struct {
	redisRepo             repos.PollingRedisRepoInterface
	brokerPollingRepo     repos.PollingBrokerRepository
	brokerCredentialsRepo repos.CredentialBrokerRepository
	httpRepo              repos.PollingHTTPRepository
	credentialHTTP        repos.CredentialHTTPRepository
}

func NewPollingService(repoRedis repos.PollingRedisRepoInterface, actionBroker repos.PollingBrokerRepository, repoHTTP repos.PollingHTTPRepository, credentialRepo repos.CredentialHTTPRepository, credentialBroker repos.CredentialBrokerRepository) repos.PollingService {
	return &PollingServiceImpl{
		redisRepo:             repoRedis,
		brokerPollingRepo:     actionBroker,
		brokerCredentialsRepo: credentialBroker,
		httpRepo:              repoHTTP,
		credentialHTTP:        credentialRepo,
	}
}

func (a *PollingServiceImpl) GetContentGoogleSheetByID(actionID *string, userID *string) (data *string, err error) {
	if actionID == nil || userID == nil {
		return nil, fmt.Errorf("cannot by empty")
	}

	data, err = a.getAllContentFromGoogleSheets(actionID, userID)
	log.Printf("%s", fmt.Sprintf("%v", *data))
	return data, err
}
