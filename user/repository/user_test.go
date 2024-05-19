package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tz3/go-simple/user/model"
)

// Test_CreateUser(t *testing.T) -> TBD

func Test_GetUserByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	repo := &UserRepositoryImpl{
		DB: mockDB,
	}

	testCases := []struct {
		name           string
		userID         int
		expectedUser   *model.User
		expectedErr    error
		expectedRows   *sqlmock.Rows
		expectedCalled int
	}{
		{
			name:   "Success_UserFound",
			userID: 123,
			expectedUser: &model.User{
				ID:        123,
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedErr: nil,
			expectedRows: sqlmock.NewRows([]string{"first_name", "last_name", "id"}).
				AddRow("John", "Doe", 123),
			expectedCalled: 1,
		},
		{
			name:         "UserNotFound",
			userID:       456,
			expectedUser: nil,
			expectedErr:  ErrUserNotFound,
			expectedRows: sqlmock.NewRows([]string{"first_name", "last_name", "id"}).
				AddRow("Ahmed", "Khalid", 123),
			expectedCalled: 1,
		},
		{
			name:           "Error_OtherError",
			userID:         789,
			expectedUser:   nil,
			expectedErr:    errors.New("some error"),
			expectedRows:   sqlmock.NewRows([]string{"first_name", "last_name", "id"}).AddRow("Weak", "Halim", 123),
			expectedCalled: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT first_name, last_name, id FROM users WHERE id = ?").
				WithArgs(tc.userID).
				WillReturnRows(tc.expectedRows).
				WillReturnError(tc.expectedErr)

			user, err := repo.GetUserByID(context.Background(), tc.userID)

			if !reflect.DeepEqual(user, tc.expectedUser) {
				t.Errorf("expected user %+v, got %+v", tc.expectedUser, user)
			}

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %+v, got %+v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}
