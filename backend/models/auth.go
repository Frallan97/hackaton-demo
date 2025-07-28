package models

// LoginInput represents the input for login requests.
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response for login requests.
type LoginResponse struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}
