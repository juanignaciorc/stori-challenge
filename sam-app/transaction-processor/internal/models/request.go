package models

// RequestBody represents the expected structure of the POST request body
type RequestBody struct {
	Email string `json:"email" validate:"required,email"`
}