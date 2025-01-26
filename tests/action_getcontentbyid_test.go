package tests

import (
	// "errors"
	"errors"
	"polling_websocket/pkg/domain/models"
	"polling_websocket/pkg/domain/repos"
	"polling_websocket/pkg/domain/services"
	"testing"

	reposmocks "polling_websocket/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPollingServiceImpl_GetContentActionByID(t *testing.T) {
	validActionID := "valid-action-id"
	validUserID := "valid-user-id"
	// validURL := "http://valid-url.com"
	validContent := "content-data"

	type fields struct {
		redisRepo             repos.PollingRedisRepoInterface
		brokerPollingRepo     repos.PollingBrokerRepository
		brokerCredentialsRepo repos.CredentialBrokerRepository
		httpRepo              repos.PollingHTTPRepository
		credentialHTTP        repos.CredentialHTTPRepository
	}

	type args struct {
		actionID *string
		userID   *string
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantData   *string
		wantErr    bool
		setupMocks func(*fields)
	}{
		{
			name:     "Happy Path - Valid action and user",
			args:     args{actionID: &validActionID, userID: &validUserID},
			wantData: &validContent,
			wantErr:  false,
			setupMocks: func(f *fields) {
				httpMock := new(reposmocks.PollingHTTPRepository)
				httpMock.On(
					"GetActionByID",
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("uint64"),
				).Return(&models.ResponsePollingActionID{
					Data: []models.RequestGoogleAction{
						{Data: validContent},
					},
				}, nil)

				f.httpRepo = httpMock
			},
		},
		{
			name:     "HTTP Request Failure",
			args:     args{actionID: &validActionID, userID: &validUserID},
			wantData: nil,
			wantErr:  true,
			setupMocks: func(f *fields) {
				redisMock := new(reposmocks.PollingRedisRepoInterface)
				credMock := new(reposmocks.CredentialBrokerRepository)
				httpMock := new(reposmocks.PollingHTTPRepository)
				httpMock.On(
					"GetActionByID",
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("uint64"),
				).Return((*models.ResponsePollingActionID)(nil), errors.New("http error"))

				f.redisRepo = redisMock
				f.brokerCredentialsRepo = credMock
				f.httpRepo = httpMock
			},
		},
		{
			name:       "Nil actionID",
			args:       args{actionID: nil, userID: &validUserID},
			wantData:   nil,
			wantErr:    true,
			setupMocks: func(f *fields) {},
		},
		{
			name:       "Nil userID",
			args:       args{actionID: &validActionID, userID: nil},
			wantData:   nil,
			wantErr:    true,
			setupMocks: func(f *fields) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := fields{}
			if tt.setupMocks != nil {
				tt.setupMocks(&fields)
			}
			service := services.NewPollingService(
				fields.redisRepo,
				fields.brokerPollingRepo,
				fields.httpRepo,
				fields.credentialHTTP,
				fields.brokerCredentialsRepo,
			)

			gotData, err := service.GetContentActionByID(tt.args.actionID, tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantData, gotData)
			if mock, ok := fields.brokerCredentialsRepo.(*reposmocks.CredentialBrokerRepository); ok {
				mock.AssertExpectations(t)
			}
			if mock, ok := fields.httpRepo.(*reposmocks.PollingHTTPRepository); ok {
				mock.AssertExpectations(t)
			}
		})
	}
}
