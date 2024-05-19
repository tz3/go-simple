package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"github.com/tz3/go-simple/config"
	"github.com/tz3/go-simple/shared"
	"github.com/tz3/go-simple/user/model"
	"github.com/tz3/go-simple/user/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// Test_CreateHandler -> TBD

func Test_GetByIdHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Create a mock UserRepo
	mockUserRepo := repository.NewMockUserRepository(ctrl)

	// Create a new UserHandler instance with the mock UserRepo
	userHandler := NewUserHandler(log, mockUserRepo)

	// Define test cases with different user IDs
	testCases := []struct {
		userID            string
		expectedStatus    int
		expectedUser      *model.User
		expectedErr       error
		expectedErrCode   string
		expectedRepoCalls int
	}{
		{
			userID:            "123",
			expectedStatus:    http.StatusOK,
			expectedUser:      &model.User{ID: 123, FirstName: "John", LastName: "Doe"},
			expectedErr:       nil,
			expectedErrCode:   "",
			expectedRepoCalls: 1,
		},
		{
			userID:            "456",
			expectedStatus:    http.StatusNotFound,
			expectedUser:      nil,
			expectedErr:       repository.ErrUserNotFound,
			expectedErrCode:   config.ErrUserNotFound,
			expectedRepoCalls: 1,
		},
		{
			userID:            "",
			expectedStatus:    http.StatusBadRequest,
			expectedUser:      nil,
			expectedErr:       repository.ErrUserNotFound,
			expectedErrCode:   config.ErrInvalidUserID,
			expectedRepoCalls: 0,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		// Create a request with the specified user ID
		req := httptest.NewRequest(http.MethodGet, "/users?userID="+tc.userID, nil)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Expect GetUserByID to be called with the corresponding user ID and return the expected user and error
		id, _ := strconv.Atoi(tc.userID)
		mockUserRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Eq(id)).Return(tc.expectedUser, tc.expectedErr).Times(tc.expectedRepoCalls)

		// Call the GetByIdHandler with the request and response recorder
		userHandler.GetByIdHandler(rr, req)

		// Check the response status code
		if rr.Code != tc.expectedStatus {
			t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
		}

		// Assert the response body content
		if tc.expectedUser != nil {
			var responseUser model.User
			err = json.NewDecoder(rr.Body).Decode(&responseUser)
			if err != nil {
				t.Errorf("failed to decode response body: %v", err)
			}

			// Check if the returned user matches the expected user
			if !reflect.DeepEqual(responseUser, *tc.expectedUser) {
				t.Errorf("expected response user %+v, got %+v", *tc.expectedUser, responseUser)
			}
		}

		// Assert the error response code and message
		if tc.expectedErr != nil {
			var responseError shared.ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&responseError)
			if err != nil {
				t.Errorf("failed to decode error response body: %v", err)
			}
		}
	}
}
