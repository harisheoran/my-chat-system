package main

type AuthRequestBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // Use `any` for generic data
}

type NotFoundResponse struct {
	Message string
}

type NotAuthorizedResponse struct {
	Message string
}
