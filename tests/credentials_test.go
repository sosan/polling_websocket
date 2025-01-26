package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	// "polling_websocket/mocks"
	"polling_websocket/mocks"
	"polling_websocket/pkg/domain/models"
	"polling_websocket/pkg/infra/brokerclient"
	"polling_websocket/pkg/infra/httpclient"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClickhouseConfig struct {
	URI    string
	Token  string
	EnvMap map[string]string
}

func (m mockClickhouseConfig) GetEnv(key, fallback string) string {
	if val, ok := m.EnvMap[key]; ok {
		return val
	}
	return fallback
}

func (m mockClickhouseConfig) GetClickhouseURI() string {
	return m.URI
}

func (m mockClickhouseConfig) GetClickhouseToken() string {
	return m.Token
}

func TestCredentialHTTPRepository_GetCredentialByID(t *testing.T) {
	type fields struct {
		DatabaseHTTPURL string
		Token           string
		Client          httpclient.HTTPClient
	}

	type args struct {
		userID       *string
		credentialID *string
		limitCount   uint64
	}

	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.RequestExchangeCredential
		wantErr bool
	}{
		{
			name: "successful request with valid credential",
			fields: fields{
				DatabaseHTTPURL: "http://example.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						expectedPath := "/credential_id_data.json"
						if req.URL.Path != expectedPath {
							return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
						}

						expectedParams := url.Values{
							"token":         []string{"test-token"},
							"credential_id": []string{"cred123"},
							"user_id":       []string{"user456"},
							"limit_count":   []string{"10"},
						}

						if !reflect.DeepEqual(req.URL.Query(), expectedParams) {
							return nil, fmt.Errorf("unexpected query params\ngot: %v\nwant: %v", req.URL.Query(), expectedParams)
						}

						jsonBody := `{
							"data": [{
								"id": "cred123",
								"name": "Test Credential",
								"created_at": "2023-01-01T12:00:00Z",
								"data": {
									"clientId": "client123",
									"clientSecret": "secret123",
									"redirectURL": "http://example.com"
								}
							}],
							"rows": 1,
							"statistics": {
								"elapsed": 0.1,
								"rows_read": 1,
								"bytes_read": 100
							}
						}`

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
				userID:       stringPtr("user456"),
				credentialID: stringPtr("cred123"),
				limitCount:   10,
			},
			want: &models.RequestExchangeCredential{
				ID:   "cred123",
				Name: "Test Credential",
				Data: models.DataCredential{
					ClientID:     "client123",
					ClientSecret: "secret123",
					RedirectURL:  "http://example.com",
				},
				CreatedAt: customTime(testTime),
			},
			wantErr: false,
		},
		{
			name: "invalid URL format",
			fields: fields{
				DatabaseHTTPURL: "://invalid-url",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:       stringPtr("user456"),
				credentialID: stringPtr("cred123"),
				limitCount:   10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "HTTP client error",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("connection timeout")
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:       stringPtr("user456"),
				credentialID: stringPtr("cred123"),
				limitCount:   10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-200 status code",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("server error")),
						}, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:       stringPtr("user456"),
				credentialID: stringPtr("cred123"),
				limitCount:   10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid JSON response",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("{invalid json")),
						}, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:       stringPtr("user456"),
				credentialID: stringPtr("cred123"),
				limitCount:   10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple credentials found",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						jsonBody := `{
							"data": [
								{"id": "cred1"},
								{"id": "cred2"}
							],
							"rows": 2
						}`

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
				userID:       stringPtr("user456"),
				credentialID: stringPtr("cred123"),
				limitCount:   10,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configuración del repositorio
			repo := httpclient.NewCredentialRepository(
				tt.fields.Client,
				mockClickhouseConfig{
					URI:   tt.fields.DatabaseHTTPURL,
					Token: tt.fields.Token,
				},
			)

			// Ejecución del método bajo test
			got, err := repo.GetCredentialByID(tt.args.userID, tt.args.credentialID, tt.args.limitCount)

			// Validación de errores
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentialByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validación de resultados
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentialByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentialHTTPRepository_GetAllCredentials(t *testing.T) {
	// Helper functions
	strPtr := func(s string) *string { return &s }
	customTime := func(t time.Time) *models.CustomTime {
		return &models.CustomTime{Time: t}
	}

	type fields struct {
		DatabaseHTTPURL string
		Token           string
		Client          httpclient.HTTPClient
	}

	type args struct {
		userID     *string
		limitCount uint64
	}

	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]models.RequestExchangeCredential
		wantErr bool
	}{
		{
			name: "successful request with multiple credentials",
			fields: fields{
				DatabaseHTTPURL: "http://example.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						expectedPath := "/all_credentials_data.json"
						if req.URL.Path != expectedPath {
							return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
						}

						expectedParams := url.Values{
							"token":       []string{"test-token"},
							"user_id":     []string{"user456"},
							"limit_count": []string{"10"},
						}

						if !reflect.DeepEqual(req.URL.Query(), expectedParams) {
							return nil, fmt.Errorf("unexpected query params\ngot: %v\nwant: %v", req.URL.Query(), expectedParams)
						}

						jsonBody := `{
							"data": [
								{
									"id": "cred1",
									"name": "Credential 1",
									"created_at": "2023-01-01T12:00:00Z",
									"data": {
										"clientId": "client1",
										"clientSecret": "secret1",
										"redirectURL": "http://example.com/1"
									}
								},
								{
									"id": "cred2",
									"name": "Credential 2",
									"created_at": "2023-01-01T12:00:00Z",
									"data": {
										"clientId": "client2",
										"clientSecret": "secret2",
										"redirectURL": "http://example.com/2"
									}
								}
							],
							"rows": 2
						}`

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
				userID:     strPtr("user456"),
				limitCount: 10,
			},
			want: &[]models.RequestExchangeCredential{
				{
					ID:   "cred1",
					Name: "Credential 1",
					Data: models.DataCredential{
						ClientID:     "client1",
						ClientSecret: "secret1",
						RedirectURL:  "http://example.com/1",
					},
					CreatedAt: customTime(testTime),
				},
				{
					ID:   "cred2",
					Name: "Credential 2",
					Data: models.DataCredential{
						ClientID:     "client2",
						ClientSecret: "secret2",
						RedirectURL:  "http://example.com/2",
					},
					CreatedAt: customTime(testTime),
				},
			},
			wantErr: false,
		},
		{
			name: "empty credentials list",
			fields: fields{
				DatabaseHTTPURL: "http://example.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						jsonBody := `{
							"data": [],
							"rows": 0
						}`
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
				userID:     strPtr("user456"),
				limitCount: 10,
			},
			want:    &[]models.RequestExchangeCredential{},
			wantErr: false,
		},
		{
			name: "invalid URL format",
			fields: fields{
				DatabaseHTTPURL: "://invalid-url",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:     strPtr("user456"),
				limitCount: 10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "HTTP client error",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("connection timeout")
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:     strPtr("user456"),
				limitCount: 10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-200 status code",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("server error")),
						}, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:     strPtr("user456"),
				limitCount: 10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid JSON response",
			fields: fields{
				DatabaseHTTPURL: "http://valid-url.com",
				Token:           "test-token",
				Client: &MockHTTPClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("{invalid json")),
						}, nil
					},
					DoRequestFunc: func(_, _, _ string, _ interface{}) ([]byte, error) {
						return nil, nil
					},
				},
			},
			args: args{
				userID:     strPtr("user456"),
				limitCount: 10,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := httpclient.NewCredentialRepository(
				tt.fields.Client,
				mockClickhouseConfig{
					URI:   tt.fields.DatabaseHTTPURL,
					Token: tt.fields.Token,
				},
			)

			got, err := repo.GetAllCredentials(tt.args.userID, tt.args.limitCount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockKafkaClient struct {
	produceFunc func(topic string, key, value []byte) error
	calls       int
}

func (m *MockKafkaClient) Close() {
}

func (m *MockKafkaClient) Produce(topic string, key, value []byte) error {
	m.calls++
	return m.produceFunc(topic, key, value)
}

func TestCredentialKafkaRepository_UpdateCredential(t *testing.T) {
	createCredential := func() *models.RequestExchangeCredential {
		return &models.RequestExchangeCredential{
			ID:   "cred123",
			Name: "Test Credential",
			Data: models.DataCredential{
				ClientID:     "client123",
				ClientSecret: "secret123",
			},
		}
	}

	tests := []struct {
		name       string
		credential *models.RequestExchangeCredential
		setupMock  func(*mocks.KafkaClient)
		wantSended bool
		wantCalls  int
	}{
		{
			name:       "successful update",
			credential: createCredential(),
			setupMock: func(m *mocks.KafkaClient) {
				m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.AnythingOfType("[]uint8")).
					Return(nil).Once()
			},
			wantSended: true,
			wantCalls:  1,
		},
		{
			name:       "retry success after failures",
			credential: createCredential(),
			setupMock: func(m *mocks.KafkaClient) {
				m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.Anything).
					Return(errors.New("connection error")).Twice()
				m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.Anything).
					Return(nil).Once()
			},
			wantSended: true,
			wantCalls:  3,
		},
		// {
		// 	name:       "max retry attempts exceeded", // Max retries 10 expect waiting time more than 2 minutes
		// 	credential: createCredential(),
		// 	setupMock: func(m *mocks.KafkaClient) {
		// 		m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.Anything).
		// 			Return(errors.New("connection error")).Times(models.MaxAttempts)
		// 	},
		// 	wantSended: false,
		// 	wantCalls:  models.MaxAttempts - 1,
		// },
		{
			name:       "invalid credential data",
			credential: nil,
			setupMock: func(m *mocks.KafkaClient) {
				// No se esperan llamadas
			},
			wantSended: false,
			wantCalls:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockKafka := new(mocks.KafkaClient)
			if tt.setupMock != nil {
				tt.setupMock(mockKafka)
			}

			repo := brokerclient.NewCredentialKafkaRepository(mockKafka)

			got := repo.UpdateCredential(tt.credential)

			assert.Equal(t, tt.wantSended, got)
			mockKafka.AssertNumberOfCalls(t, "Produce", tt.wantCalls)

			if tt.wantSended && tt.credential != nil {
				expectedCommand := brokerclient.CredentialCommand{
					Type:       models.UpdateCommand,
					Credential: tt.credential,
				}

				mockKafka.AssertCalled(t, "Produce", models.CredentialTopicName, []byte(tt.credential.ID), mock.MatchedBy(func(value []byte) bool {
					var actualCommand brokerclient.CredentialCommand
					if err := json.Unmarshal(value, &actualCommand); err != nil {
						return false
					}

					expected, _ := json.Marshal(expectedCommand)
					actual, _ := json.Marshal(actualCommand)
					return string(expected) == string(actual)
				}))
			}
		})
	}
}

func TestCredentialKafkaRepository_PublishCommand(t *testing.T) {
	createValidCommand := func() brokerclient.CredentialCommand {
		return brokerclient.CredentialCommand{
			Type: models.UpdateCommand,
			Credential: &models.RequestExchangeCredential{
				ID:   "cred123",
				Name: "Test Credential",
				Data: models.DataCredential{
					ClientID:     "client123",
					ClientSecret: "secret123",
				},
			},
		}
	}

	tests := []struct {
		name       string
		command    brokerclient.CredentialCommand
		key        string
		setupMock  func(*mocks.KafkaClient)
		wantSended bool
		wantCalls  int
	}{
		{
			name:    "successful publish",
			command: createValidCommand(),
			key:     "cred123",
			setupMock: func(m *mocks.KafkaClient) {
				m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.AnythingOfType("[]uint8")).
					Return(nil).Once()
			},
			wantSended: true,
			wantCalls:  1,
		},
		// {
		// 	name:    "retry success after failures", // more than 2 minutes
		// 	command: createValidCommand(),
		// 	key:     "cred123",
		// 	setupMock: func(m *mocks.KafkaClient) {
		// 		m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.Anything).
		// 			Return(errors.New("connection error")).Twice()
		// 		m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.Anything).
		// 			Return(nil).Once()
		// 	},
		// 	wantSended: true,
		// 	wantCalls:  3,
		// },
		// {
		// 	name:    "max retry attempts exceeded",
		// 	command: createValidCommand(),
		// 	key:     "cred123",
		// 	setupMock: func(m *mocks.KafkaClient) {
		// 		m.On("Produce", models.CredentialTopicName, []byte("cred123"), mock.Anything).
		// 			Return(errors.New("connection error")).Times(models.MaxAttempts)
		// 	},
		// 	wantSended: false,
		// 	wantCalls:  models.MaxAttempts - 1,
		// },
		{
			name:    "empty key",
			command: createValidCommand(),
			key:     "",
			setupMock: func(m *mocks.KafkaClient) {
				m.On("Produce", models.CredentialTopicName, []byte(""), mock.Anything).
					Return(nil).Once()
			},
			wantSended: true,
			wantCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockKafka := new(mocks.KafkaClient)
			if tt.setupMock != nil {
				tt.setupMock(mockKafka)
			}

			repo := brokerclient.NewCredentialKafkaRepository(mockKafka)
			got := repo.PublishCommand(tt.command, tt.key)

			assert.Equal(t, tt.wantSended, got)
			mockKafka.AssertNumberOfCalls(t, "Produce", tt.wantCalls)

			if tt.wantSended {
				mockKafka.AssertCalled(t, "Produce",
					models.CredentialTopicName,
					[]byte(tt.key),
					mock.MatchedBy(func(value []byte) bool {
						var actualCommand brokerclient.CredentialCommand
						if err := json.Unmarshal(value, &actualCommand); err != nil {
							return false
						}
						expected, _ := json.Marshal(tt.command)
						actual, _ := json.Marshal(actualCommand)
						return string(expected) == string(actual)
					}),
				)
			}
		})
	}
}
