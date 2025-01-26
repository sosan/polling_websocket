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

type PollingHTTPRepository struct {
	DatabaseHTTPURL string
	Token           string
	Client          HTTPClient
}

func NewPollingClientHTTP(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *PollingHTTPRepository {
	return &PollingHTTPRepository{
		Client:          client,
		DatabaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		Token:           clickhouseConfig.GetClickhouseToken(),
	}
}

// GetActionByID retrieves action details by action ID from the polling HTTP repository.
// It constructs a URL with the provided parameters, sends a GET request, and decodes the response.
//
// Parameters:
//   - actionID: Pointer to the action ID string.
//   - userID: Pointer to the user ID string.
//   - commandType: Type of the command as a string.
//   - limitCount: Limit count as an unsigned 64-bit integer.
//
// Returns:
//   - data: Pointer to the ResponsePollingActionID model containing the action details.
//   - err: Error if any occurred during the process.
//
// Errors:
//   - Returns an error if the URL cannot be parsed.
//   - Returns an error if the HTTP request cannot be generated.
//   - Returns an error if the HTTP request fails.
//   - Returns an error if the response status is not OK (200).
//   - Returns an error if the response body is nil or cannot be decoded.
//
// @Summary Get action details by action ID
// @Description Retrieves action details by action ID from the polling HTTP repository
// @Tags actions
// @Accept json
// @Produce json
// @Param action_id query string true "Action ID"
// @Param user_id query string true "User ID"
// @Param command_type query string true "Command Type"
// @Param limit_count query int true "Limit Count"
// @Success 200 {object} models.ResponsePollingActionID
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /action_workflow_data.json [get]
func (a *PollingHTTPRepository) GetActionByID(actionID *string, userID *string, commandType string, limitCount uint64) (data *models.ResponsePollingActionID, err error) {
	u, err := url.Parse(a.DatabaseHTTPURL + "/action_workflow_data.json")
	if err != nil {
		log.Printf("ERROR | polling httpclient cannot parse url %v", err)
		return nil, err
	}

	q := u.Query()
	q.Set("token", a.Token)
	q.Set("action_id", *actionID)
	q.Set("user_id", *userID)
	q.Set("command_type", commandType)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("ERROR | polling httpclient cannot generate request %v", err)
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Printf("ERROR | polling httpclient not response %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("ERROR | response body is nil")
	}

	var result *models.ResponsePollingActionID
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}
	return result, nil
}
