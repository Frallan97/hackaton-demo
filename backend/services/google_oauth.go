package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/frallan97/react-go-app-backend/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
)

// GoogleOAuthService handles Google OAuth operations
type GoogleOAuthService struct {
	config *oauth2.Config
}

// NewGoogleOAuthService creates a new Google OAuth service
func NewGoogleOAuthService(clientID, clientSecret, redirectURL string) *GoogleOAuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthService{
		config: config,
	}
}

// GetAuthURL returns the Google OAuth authorization URL
func (g *GoogleOAuthService) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (g *GoogleOAuthService) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}

// GetUserInfo retrieves user information from Google using the access token
func (g *GoogleOAuthService) GetUserInfo(token *oauth2.Token) (*models.GoogleUserInfo, error) {
	ctx := context.Background()
	client := g.config.Client(ctx, token)

	// Create OAuth2 service
	oauth2Service, err := googleoauth2.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth2 service: %w", err)
	}

	// Get user info
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	verifiedEmail := false
	if userInfo.VerifiedEmail != nil {
		verifiedEmail = *userInfo.VerifiedEmail
	}

	return &models.GoogleUserInfo{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		VerifiedEmail: verifiedEmail,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}, nil
}

// ValidateToken validates a Google access token
func (g *GoogleOAuthService) ValidateToken(token *oauth2.Token) error {
	ctx := context.Background()
	client := g.config.Client(ctx, token)

	// Make a request to Google's userinfo endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	return nil
}
