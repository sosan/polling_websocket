package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"polling_websocket/pkg/config"
	"polling_websocket/pkg/domain/models"
)

type CredentialHTTPRepository struct {
	client          HTTPClient
	databaseHTTPURL string
	token           string
}

func NewCredentialRepository(httpCli HTTPClient, clickhouseConfig config.ClickhouseConfig) *CredentialHTTPRepository {
	return &CredentialHTTPRepository{
		client:          httpCli,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (c *CredentialHTTPRepository) GetCredentialByID(userID *string, credentialID *string, limitCount uint64) (*models.RequestExchangeCredential, error) {
	u, err := url.Parse(c.databaseHTTPURL + "/credential_id_data.json")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", c.token)
	q.Set("credential_id", *credentialID)
	q.Set("user_id", *userID)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result *models.InfoCredentials

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}
	if len(*result.Data) > 1 {
		// length cannot be more than 1
		log.Printf("ERROR | duplicated ID: %v", result.Data)
		return nil, fmt.Errorf("ERROR | duplicated id token")
	}
	return &(*result.Data)[0], nil
}

func (c *CredentialHTTPRepository) GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error) {
	u, err := url.Parse(c.databaseHTTPURL + "/all_credentials_data.json")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", c.token)
	q.Set("user_id", *userID)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result *models.InfoCredentials

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}

	return result.Data, nil
}
