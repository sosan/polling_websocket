package repos

import "polling_websocket/pkg/domain/models"

type CredentialHTTPRepository interface {
	GetCredentialByID(actionID *string, userID *string, limitCount uint64) (*models.RequestExchangeCredential, error)
	GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error)
}
