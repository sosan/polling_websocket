package di

import (
	"polling_websocket/pkg/config"
	"polling_websocket/pkg/dimodel"
	"polling_websocket/pkg/domain/models"
	"polling_websocket/pkg/domain/services"
	"polling_websocket/pkg/infra/brokerclient"
	"polling_websocket/pkg/infra/httpclient"
	"polling_websocket/pkg/infra/redisclient"
	"polling_websocket/pkg/interfaces/controllers"
)

func InitDependencies() *dimodel.Dependencies {
	configZitadel := config.NewZitaldelEnvConfig()
	kafkaConfig := config.NewKafkaEnvConfig()
	clickhouseConfig := config.NewClickhouseEnvConfig()

	// init autentication
	authContext := controllers.NewAuthContext(configZitadel)
	authService := authContext.GetAuthService()
	authController := authContext.GetAuthController()

	credentialBrokerClient := brokerclient.NewBrokerClient(kafkaConfig)
	repoCredentialBroker := brokerclient.NewCredentialKafkaRepository(credentialBrokerClient)

	actionsHTTPClient := httpclient.NewClientImpl(models.TimeoutRequest)
	credentialHTTPClient := httpclient.NewClientImpl(models.TimeoutRequest)
	repoCredentialHTTP := httpclient.NewCredentialRepository(credentialHTTPClient, clickhouseConfig)
	actionsRedisClient := redisclient.NewRedisClient()
	actionsBrokerClient := brokerclient.NewBrokerClient(kafkaConfig)
	repoPollingRedis := redisclient.NewPollingRepository(actionsRedisClient)
	repoPollingBroker := brokerclient.NewPollingKafkaRepository(actionsBrokerClient)
	actionsRepo := httpclient.NewPollingClientHTTP(actionsHTTPClient, clickhouseConfig)
	actionsService := services.NewPollingService(repoPollingRedis, repoPollingBroker, actionsRepo, repoCredentialHTTP, repoCredentialBroker)
	actionsController := controllers.NewPollingController(actionsService)

	return &dimodel.Dependencies{
		AuthService:       &authService,
		AuthController:    authController,
		PollingController: actionsController,
	}
}
