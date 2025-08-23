package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/frallan97/hackaton-demo-backend/database"
	"github.com/frallan97/hackaton-demo-backend/events"
	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/frallan97/hackaton-demo-backend/services"
	"github.com/frallan97/hackaton-demo-backend/utils"
)

// AuthController handles authentication-related endpoints
type AuthController struct {
	dbManager          *database.DBManager
	userService        *services.UserService
	jwtService         *services.JWTService
	googleOAuthService *services.GoogleOAuthService
	eventService       *events.EventService
}

// NewAuthController creates a new auth controller
func NewAuthController(dbManager *database.DBManager, userService *services.UserService, jwtService *services.JWTService, googleOAuthService *services.GoogleOAuthService, eventService *events.EventService) *AuthController {
	return &AuthController{
		dbManager:          dbManager,
		userService:        userService,
		jwtService:         jwtService,
		googleOAuthService: googleOAuthService,
		eventService:       eventService,
	}
}

// GoogleLoginHandler handles Google OAuth login
// @Summary     Google OAuth Login
// @Description Authenticate user with Google OAuth
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       login  body   models.LoginRequest  true  "Google OAuth code"
// @Success     200   {object}  utils.APIResponse{data=models.AuthResponse}
// @Failure     400   {object}  utils.APIResponse
// @Failure     405   {object}  utils.APIResponse
// @Failure     500   {object}  utils.APIResponse
// @Router      /api/auth/google/login [post]
func (ac *AuthController) GoogleLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteBadRequest(w, "Invalid request body", err)
			return
		}

		// Validate input
		if req.Code == "" {
			utils.WriteValidationError(w, map[string]string{
				"code": "Authorization code is required",
			})
			return
		}

		// Exchange authorization code for access token
		token, err := ac.googleOAuthService.ExchangeCodeForToken(req.Code)
		if err != nil {
			utils.WriteBadRequest(w, "Failed to exchange authorization code", err)
			return
		}

		// Get user info from Google
		googleUserInfo, err := ac.googleOAuthService.GetUserInfo(token)
		if err != nil {
			utils.WriteBadRequest(w, "Failed to get user info from Google", err)
			return
		}

		// Check if user exists in our database
		user, err := ac.userService.GetUserByGoogleID(googleUserInfo.ID)
		if err != nil {
			// Log the actual error for debugging
			fmt.Printf("Database error getting user by Google ID: %v\n", err)
			utils.WriteInternalServerError(w, "Database error while retrieving user", err)
			return
		}

		if user == nil {
			// Create new user
			userData := &models.UserCreate{
				Email:    googleUserInfo.Email,
				Name:     googleUserInfo.Name,
				Picture:  googleUserInfo.Picture,
				GoogleID: googleUserInfo.ID,
			}

			user, err = ac.userService.CreateUser(userData)
			if err != nil {
				// Log the actual error for debugging
				fmt.Printf("Failed to create user: %v\n", err)
				utils.WriteInternalServerError(w, "Failed to create user account", err)
				return
			}

			// Publish user created event
			if ac.eventService != nil {
				if err := ac.eventService.PublishUserCreated(user.ID, user.Email, user.Name); err != nil {
					fmt.Printf("Warning: Failed to publish user created event: %v\n", err)
				}
			}
		} else {
			// Update last login time
			err = ac.userService.UpdateUserLastLogin(user.ID)
			if err != nil {
				// Log error but don't fail the login
				fmt.Printf("failed to update last login: %v\n", err)
			}

			// Update profile if needed
			if user.Name != googleUserInfo.Name || user.Picture != googleUserInfo.Picture {
				user, err = ac.userService.UpdateUserProfile(user.ID, googleUserInfo.Name, googleUserInfo.Picture)
				if err != nil {
					// Log error but don't fail the login
					fmt.Printf("failed to update profile: %v\n", err)
				}
			}
		}

		// Publish user login event
		if ac.eventService != nil {
			if err := ac.eventService.PublishUserLogin(user.ID, user.Email, user.Name); err != nil {
				fmt.Printf("Warning: Failed to publish user login event: %v\n", err)
			}
		}

		// Generate JWT tokens
		accessToken, refreshToken, err := ac.jwtService.GenerateTokens(user)
		if err != nil {
			utils.WriteInternalServerError(w, "Failed to generate authentication tokens", err)
			return
		}

		// Create response
		response := models.AuthResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(ac.jwtService.GetTokenExpiry().Seconds()),
		}

		utils.WriteOK(w, response, "Login successful")
	}
}

// GetAuthURLHandler returns the Google OAuth authorization URL
// @Summary     Get Google OAuth URL
// @Description Get the Google OAuth authorization URL
// @Tags        auth
// @Produce     json
// @Success     200   {object}  utils.APIResponse{data=map[string]string}
// @Router      /api/auth/google/url [get]
func (ac *AuthController) GetAuthURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		state := "random-state-string" // In production, generate a secure random state
		authURL := ac.googleOAuthService.GetAuthURL(state)

		response := map[string]string{
			"auth_url": authURL,
			"state":    state,
		}

		utils.WriteOK(w, response, "Google OAuth URL generated successfully")
	}
}

// RefreshTokenResponse represents the response for token refresh
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

// RefreshTokenHandler refreshes an access token using a refresh token
// @Summary     Refresh Token
// @Description Refresh access token using refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       refresh  body   models.RefreshTokenRequest  true  "Refresh token"
// @Success     200   {object}  utils.APIResponse{data=RefreshTokenResponse}
// @Failure     400   {object}  utils.APIResponse
// @Failure     405   {object}  utils.APIResponse
// @Router      /api/auth/refresh [post]
func (ac *AuthController) RefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		var req models.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteBadRequest(w, "Invalid request body", err)
			return
		}

		// Validate input
		if req.RefreshToken == "" {
			utils.WriteValidationError(w, map[string]string{
				"refresh_token": "Refresh token is required",
			})
			return
		}

		// Refresh the access token
		newAccessToken, err := ac.jwtService.RefreshToken(req.RefreshToken)
		if err != nil {
			utils.WriteBadRequest(w, "Invalid refresh token", err)
			return
		}

		response := &RefreshTokenResponse{
			AccessToken: newAccessToken,
			TokenType:   "Bearer",
			ExpiresIn:   fmt.Sprintf("%d", int(ac.jwtService.GetTokenExpiry().Seconds())),
		}

		utils.WriteOK(w, response, "Token refreshed successfully")
	}
}

// GetMeHandler returns the current user's information
// @Summary     Get Current User
// @Description Get current user information
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200   {object}  utils.APIResponse{data=models.User}
// @Failure     401   {object}  utils.APIResponse
// @Failure     404   {object}  utils.APIResponse
// @Failure     405   {object}  utils.APIResponse
// @Failure     500   {object}  utils.APIResponse
// @Router      /api/auth/me [get]
func (ac *AuthController) GetMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteUnauthorized(w, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.WriteUnauthorized(w, "Invalid authorization header format")
			return
		}

		// Validate token
		claims, err := ac.jwtService.ValidateToken(tokenString)
		if err != nil {
			utils.WriteUnauthorized(w, "Invalid token")
			return
		}

		// Get user from database
		user, err := ac.userService.GetUserByID(claims.UserID)
		if err != nil {
			utils.WriteInternalServerError(w, "Database error while retrieving user", err)
			return
		}

		if user == nil {
			utils.WriteNotFound(w, "User not found")
			return
		}

		utils.WriteOK(w, user, "User information retrieved successfully")
	}
}

// LogoutResponse represents the response for logout
type LogoutResponse struct {
	Message string `json:"message"`
}

// LogoutHandler handles user logout
// @Summary     Logout
// @Description Logout user (client should discard tokens)
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200   {object}  utils.APIResponse{data=LogoutResponse}
// @Failure     405   {object}  utils.APIResponse
// @Router      /api/auth/logout [post]
func (ac *AuthController) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		// Extract token from Authorization header to get user info for event publishing
		authHeader := r.Header.Get("Authorization")
		var userID int
		var userEmail string
		var userName string

		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString != authHeader {
				// Validate token to get user info
				if claims, err := ac.jwtService.ValidateToken(tokenString); err == nil {
					userID = claims.UserID
					// Get user details for event
					if user, err := ac.userService.GetUserByID(userID); err == nil && user != nil {
						userEmail = user.Email
						userName = user.Name
					}
				}
			}
		}

		// Publish logout event if we have user info
		if ac.eventService != nil && userID > 0 {
			if err := ac.eventService.PublishUserLogout(userID, userEmail, userName); err != nil {
				fmt.Printf("Warning: Failed to publish user logout event: %v\n", err)
			}
		}

		// In a stateless JWT system, logout is handled client-side
		// The server can't invalidate JWT tokens, so we just return success
		// For enhanced security, you could implement a token blacklist using Redis

		response := &LogoutResponse{
			Message: "Logged out successfully",
		}

		utils.WriteOK(w, response, "Logout successful")
	}
}
