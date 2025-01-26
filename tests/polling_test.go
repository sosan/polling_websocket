package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"polling_websocket/pkg/domain/models"
	"polling_websocket/pkg/infra/httpclient"
	"reflect"
	"testing"
)

func TestPollingHTTPRepository_GetActionByID(t *testing.T) {
	type fields struct {
		databaseHTTPURL string
		token           string
		client          httpclient.HTTPClient
	}
	type args struct {
		actionID    *string
		userID      *string
		commandType string
		limitCount  uint64
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData *models.ResponsePollingActionID
		wantErr  bool
	}{
		{
			name: "successful request",
			fields: fields{
				databaseHTTPURL: "http://example.com",
				token:           "test-token",
				client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						if req.URL.Path != "/action_workflow_data.json" {
							return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
						}

						expectedParams := url.Values{
							"token":        []string{"test-token"},
							"action_id":    []string{"action123"},
							"user_id":      []string{"user456"},
							"command_type": []string{"type1"},
							"limit_count":  []string{"10"},
						}

						actualParams := req.URL.Query()
						if !reflect.DeepEqual(actualParams, expectedParams) {
							return nil, fmt.Errorf("unexpected query params\ngot: %v\nwant: %v", actualParams, expectedParams)
						}

						jsonBody := `{ "data": [], "rows": 0, "rows_before_limit_at_least": 0 }`
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(jsonBody)),
						}, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				actionID:    stringPtr("action123"),
				userID:      stringPtr("user456"),
				commandType: "type1",
				limitCount:  10,
			},
			wantData: &models.ResponsePollingActionID{
				Data:                   []models.RequestGoogleAction{},
				Rows:                   int64Ptr(0),
				RowsBeforeLimitAtLeast: int64Ptr(0),
				Statistics:             nil,
				Meta:                   nil,
			},
			wantErr: false,
		},
		{
			name: "invalid database URL",
			fields: fields{
				databaseHTTPURL: "://invalid-url",
				token:           "test-token",
				client:          &MockHTTPClient{},
			},
			args: args{
				actionID:    stringPtr("action123"),
				userID:      stringPtr("user456"),
				commandType: "type1",
				limitCount:  10,
			},
			wantData: nil,
			wantErr:  true,
		},
		{
			name: "HTTP client error",
			fields: fields{
				databaseHTTPURL: "http://valid-url.com",
				token:           "test-token",
				client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("connection error")
					},
				},
			},
			args: args{
				actionID:    stringPtr("action123"),
				userID:      stringPtr("user456"),
				commandType: "type1",
				limitCount:  10,
			},
			wantData: nil,
			wantErr:  true,
		},
		{
			name: "non-200 status code",
			fields: fields{
				databaseHTTPURL: "http://valid-url.com",
				token:           "test-token",
				client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("server error")),
						}, nil
					},
				},
			},
			args: args{
				actionID:    stringPtr("action123"),
				userID:      stringPtr("user456"),
				commandType: "type1",
				limitCount:  10,
			},
			wantData: nil,
			wantErr:  true,
		},
		{
			name: "invalid JSON response",
			fields: fields{
				databaseHTTPURL: "http://valid-url.com",
				token:           "test-token",
				client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("{invalid json")),
						}, nil
					},
				},
			},
			args: args{
				actionID:    stringPtr("action123"),
				userID:      stringPtr("user456"),
				commandType: "type1",
				limitCount:  10,
			},
			wantData: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &httpclient.PollingHTTPRepository{
				DatabaseHTTPURL: tt.fields.databaseHTTPURL,
				Token:           tt.fields.token,
				Client:          tt.fields.client,
			}

			gotData, err := a.GetActionByID(tt.args.actionID, tt.args.userID, tt.args.commandType, tt.args.limitCount)

			if (err != nil) != tt.wantErr {
				t.Errorf("PollingHTTPRepository.GetActionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("PollingHTTPRepository.GetActionByID() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
