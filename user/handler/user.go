package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/tz3/go-simple/config"
	"github.com/tz3/go-simple/shared"
	"github.com/tz3/go-simple/user/model"
	"github.com/tz3/go-simple/user/repository"
)

// UserHandler represents the handler for users
type UserHandler struct {
	Log      *zap.Logger
	UserRepo repository.UserRepository
}

// NewUserHandler creates a new instance of UserHandler with the given logger and user repository.
// It initializes the UserHandler struct and returns a pointer to it.
func NewUserHandler(logger *zap.Logger, userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		Log:      logger,
		UserRepo: userRepo,
	}
}

// UserHandler handles the HTTP requests for user-related operations.
// It delegates the request to the appropriate handler based on the request method.
// It logs any errors related to invalid request methods.
func (u *UserHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.GetByIdHandler(w, r)
	case http.MethodPost:
		u.CreateHandler(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		u.Log.Error("Invalid method", zap.String("method", r.Method))
	}
}

// GetByIdHandler handles the GET request to retrieve a user by ID.
// /users?UserID={id}
func (u *UserHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.FormValue("userID") // TODO:- TBD it is better if it is not query like => /users/{id}

	if idStr == "" {
		shared.HandleJSONError(w, config.ErrEmptyUserID, http.StatusBadRequest, u.Log, "Some error occurred")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.HandleJSONError(w, config.ErrInvalidUserID, http.StatusBadRequest, u.Log, "Some error occurred")
		return
	}

	user, err := u.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			shared.HandleJSONError(w, config.ErrUserNotFound, http.StatusNotFound, u.Log, "Some error occurred")
			return
		}
		shared.HandleJSONError(w, config.ErrFetchingUser, http.StatusInternalServerError, u.Log, "Some error occurred")
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		shared.HandleJSONError(w, config.ErrMarshallingUser, http.StatusInternalServerError, u.Log, "Some error occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		u.Log.Error("Error writing response", zap.Error(err))
	}
}

// CreateHandler handles the POST request to create a user.
// /users with {{user}} as payload json object
func (u *UserHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.DefaultTimeout)
	defer cancel()

	var req model.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		shared.HandleJSONError(w, "Bad Request", http.StatusBadRequest, u.Log, "Failed to decode request body")
		return
	}

	// Validate the request input
	if err := validateCreateUserRequest(&req); err != nil {
		shared.HandleJSONError(w, "Bad Request", http.StatusBadRequest, u.Log, "Validation error:", err.Error())
		return
	}

	// Map request struct to model struct
	user := mapCreateUserRequestToUser(&req)

	// Create the user using the UserRepository
	err = u.UserRepo.CreateUser(ctx, &user)
	if err != nil {
		shared.HandleJSONError(w, "Internal Server Error", http.StatusInternalServerError, u.Log, "Failed to create user")
		return
	}

	// Map model struct to response struct
	res := mapUserToCreateUserResponse(&user)

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// validateCreateUserRequest is a helper function to validate the create user request
func validateCreateUserRequest(req *model.CreateUserRequest) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		// Iterate over validation errors and concatenate them into a single error message
		var errMsg string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += err.Field() + " is a required field; "
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}

// mapCreateUserRequestToUser is a helper function to map the create user request to the user model
func mapCreateUserRequestToUser(req *model.CreateUserRequest) model.User {
	return model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}

// mapUserToCreateUserResponse is a helper function to map the user model to the create user response
func mapUserToCreateUserResponse(user *model.User) model.CreateUserResponse {
	return model.CreateUserResponse{
		User: *user,
	}
}
