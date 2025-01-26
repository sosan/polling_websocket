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
	Client          HTTPClient
	DatabaseHTTPURL string
	Token           string
}

func NewCredentialRepository(httpCli HTTPClient, clickhouseConfig config.ClickhouseConfig) *CredentialHTTPRepository {
	return &CredentialHTTPRepository{
		Client:          httpCli,
		DatabaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		Token:           clickhouseConfig.GetClickhouseToken(),
	}
}

// GetCredentialByID retrieves a credential by its ID from the HTTP repository.
// It constructs a URL with query parameters including token, credential ID, user ID, and limit count,
// then sends a GET request to the constructed URL. The response is expected to be in JSON format.
//
// Parameters:
//   - userID: A pointer to a string representing the user ID.
//   - credentialID: A pointer to a string representing the credential ID.
//   - limitCount: A uint64 representing the limit count.
//
// Returns:
//   - A pointer to models.RequestExchangeCredential containing the credential data if successful.
//   - An error if there is any issue with the request, response, or data decoding.
//
// Errors:
//   - Returns an error if the URL parsing fails.
//   - Returns an error if the HTTP request creation fails.
//   - Returns an error if the HTTP request execution fails.
//   - Returns an error if the response status code is not 200 OK.
//   - Returns an error if the response body cannot be decoded into the expected structure.
//   - Returns an error if the response contains more than one credential data entry.
//
// @Summary Get credential by ID
// @Description Retrieves a credential by its ID from the HTTP repository
// @Tags credentials
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param credential_id query string true "Credential ID"
// @Param limit_count query int true "Limit Count"
// @Success 200 {object} models.RequestExchangeCredential
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /credential_id_data.json [get]
func (c *CredentialHTTPRepository) GetCredentialByID(userID *string, credentialID *string, limitCount uint64) (*models.RequestExchangeCredential, error) {
	u, err := url.Parse(c.DatabaseHTTPURL + "/credential_id_data.json")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", c.Token)
	q.Set("credential_id", *credentialID)
	q.Set("user_id", *userID)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
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

// GetAllCredentials retrieves all credentials for a given user ID with a specified limit count.
// It sends a GET request to the configured database HTTP URL with the provided token, user ID, and limit count as query parameters.
// If the request is successful, it decodes the response body into a slice of RequestExchangeCredential models.
// 
// Parameters:
//   - userID: A pointer to a string representing the user ID for which to retrieve credentials.
//   - limitCount: An unsigned integer representing the maximum number of credentials to retrieve.
//
// Returns:
//   - A pointer to a slice of RequestExchangeCredential models containing the retrieved credentials.
//   - An error if the request fails or the response cannot be decoded.
//
// @Summary Get credential by ID
// @Description Retrieves a credential by its ID from the HTTP repository
// @Tags credentials
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param credential_id query string true "Credential ID"
// @Param limit_count query int true "Limit Count"
// @Success 200 {object} models.RequestExchangeCredential
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /credential_id_data.json [get]
func (c *CredentialHTTPRepository) GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error) {
	u, err := url.Parse(c.DatabaseHTTPURL + "/all_credentials_data.json")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", c.Token)
	q.Set("user_id", *userID)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
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
