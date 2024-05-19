package model

// User represents a user struct with ID, FirstName, and LastName
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// request

// CreateUserRequest represents the request structure for creating a user.
type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// response

// CreateUserResponse represents the response structure for creating a user.
type CreateUserResponse struct {
	User
}
