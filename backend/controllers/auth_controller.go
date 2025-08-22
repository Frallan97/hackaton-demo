package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/frallan97/hackaton-demo-backend/database"
	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/frallan97/hackaton-demo-backend/services"
)

// AuthController handles authentication-related endpoints
type AuthController struct {
	dbManager          *database.DBManager
	userService        *services.UserService
	jwtService         *services.JWTService
	googleOAuthService *services.GoogleOAuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(dbManager *database.DBManager, userService *services.UserService, jwtService *services.JWTService, googleOAuthService *services.GoogleOAuthService) *AuthController {
	return &AuthController{
		dbManager:          dbManager,
		userService:        userService,
		jwtService:         jwtService,
		googleOAuthService: googleOAuthService,
	}
}

// GoogleLoginHandler handles Google OAuth login
// @Summary     Google OAuth Login
// @Description Authenticate user with Google OAuth
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       login  body   models.LoginRequest  true  "Google OAuth code"
// @Success     200   {object}  models.AuthResponse
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /api/auth/google/login [post]
func (ac *AuthController) GoogleLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", 405)
			return
		}

		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", 400)
			return
		}

		// Exchange authorization code for access token
		token, err := ac.googleOAuthService.ExchangeCodeForToken(req.Code)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to exchange code: %v", err), 400)
			return
		}

		// Get user info from Google
		googleUserInfo, err := ac.googleOAuthService.GetUserInfo(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get user info: %v", err), 400)
			return
		}

		// Check if user exists in our database
		user, err := ac.userService.GetUserByGoogleID(googleUserInfo.ID)
		if err != nil {
			// Log the actual error for debugging
			fmt.Printf("Database error getting user by Google ID: %v\n", err)
			http.Error(w, fmt.Sprintf("database error: %v", err), 500)
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
				http.Error(w, fmt.Sprintf("failed to create user: %v", err), 500)
				return
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

		// Generate JWT tokens
		accessToken, refreshToken, err := ac.jwtService.GenerateTokens(user)
		if err != nil {
			http.Error(w, "failed to generate tokens", 500)
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetAuthURLHandler returns the Google OAuth authorization URL
// @Summary     Get Google OAuth URL
// @Description Get the Google OAuth authorization URL
// @Tags        auth
// @Produce     json
// @Success     200   {object}  map[string]string
// @Router      /api/auth/google/url [get]
func (ac *AuthController) GetAuthURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", 405)
			return
		}

		state := "random-state-string" // In production, generate a secure random state
		authURL := ac.googleOAuthService.GetAuthURL(state)

		response := map[string]string{
			"auth_url": authURL,
			"state":    state,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// RefreshTokenHandler refreshes an access token using a refresh token
// @Summary     Refresh Token
// @Description Refresh access token using refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       refresh  body   models.RefreshTokenRequest  true  "Refresh token"
// @Success     200   {object}  map[string]string
// @Failure     400   {object}  map[string]string
// @Router      /api/auth/refresh [post]
func (ac *AuthController) RefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", 405)
			return
		}

		var req models.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", 400)
			return
		}

		// Refresh the access token
		newAccessToken, err := ac.jwtService.RefreshToken(req.RefreshToken)
		if err != nil {
			http.Error(w, "invalid refresh token", 400)
			return
		}

		response := map[string]string{
			"access_token": newAccessToken,
			"token_type":   "Bearer",
			"expires_in":   fmt.Sprintf("%d", int(ac.jwtService.GetTokenExpiry().Seconds())),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetMeHandler returns the current user's information
// @Summary     Get Current User
// @Description Get current user information
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200   {object}  models.User
// @Failure     401   {object}  map[string]string
// @Router      /api/auth/me [get]
func (ac *AuthController) GetMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", 405)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", 401)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "invalid authorization header format", 401)
			return
		}

		// Validate token
		claims, err := ac.jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		// Get user from database
		user, err := ac.userService.GetUserByID(claims.UserID)
		if err != nil {
			http.Error(w, "database error", 500)
			return
		}

		if user == nil {
			http.Error(w, "user not found", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// LogoutHandler handles user logout
// @Summary     Logout
// @Description Logout user (client should discard tokens)
// @Tags        auth
// @Produce     json
// @Success     200   {object}  map[string]string
// @Router      /api/auth/logout [post]
func (ac *AuthController) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", 405)
			return
		}

		// In a stateless JWT system, logout is handled client-side
		// The server can't invalidate JWT tokens, so we just return success
		// For enhanced security, you could implement a token blacklist using Redis

		response := map[string]string{
			"message": "Logged out successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
