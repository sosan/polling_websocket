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
	databaseHTTPURL string
	token           string
	client          HTTPClient
}

func NewPollingClientHTTP(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *PollingHTTPRepository {
	return &PollingHTTPRepository{
		client:          client,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (a *PollingHTTPRepository) GetActionByID(actionID *string, userID *string, commandType string, limitCount uint64) (data *models.ResponsePollingActionID, err error) {
	u, err := url.Parse(a.databaseHTTPURL + "/action_workflow_data.json")
	if err != nil {
		log.Printf("ERROR | polling httpclient cannot parse url %v", err)
		return nil, err
	}

	q := u.Query()
	q.Set("token", a.token)
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

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("ERROR | polling httpclient not response %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result *models.ResponsePollingActionID
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}

	return result, nil
}
