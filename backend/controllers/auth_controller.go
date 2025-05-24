package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"resume.in/backend/config"
	"resume.in/backend/models"
	"resume.in/backend/utils"
)

// AuthController handles authentication endpoints
type AuthController struct {
	oauthConfig *oauth2.Config
	userRepo    models.UserRepository
	jwtSecret   string
}

// NewAuthController creates a new auth controller
func NewAuthController(cfg *config.Config, userRepo models.UserRepository) *AuthController {
	return &AuthController{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		userRepo:  userRepo,
		jwtSecret: cfg.JWTSecret,
	}
}

// GoogleUserInfo represents the user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string             `json:"token"`
	RefreshToken string             `json:"refresh_token"`
	User         models.User        `json:"user"`
	ExpiresIn    int64              `json:"expires_in"`
}

// GoogleLogin initiates the Google OAuth flow
// @Summary Initiate Google OAuth login
// @Description Redirects to Google OAuth consent page
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/google/login [get]
func (a *AuthController) GoogleLogin(c *gin.Context) {
	state := utils.GenerateRandomString(32)
	
	// Store state in session or temporary storage for verification
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	
	url := a.oauthConfig.AuthCodeURL(state)
	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
	})
}

// GoogleCallback handles the Google OAuth callback
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth callback and returns JWT tokens
// @Tags auth
// @Param code query string true "Authorization code"
// @Param state query string true "OAuth state"
// @Produce json
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/google/callback [get]
func (a *AuthController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	
	// Verify state
	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	
	// Exchange code for token
	token, err := a.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.Error("Failed to exchange token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	
	// Get user info from Google
	client := a.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.Error("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()
	
	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		utils.Error("Failed to decode user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}
	
	// Check if user exists or create new user
	user, err := a.userRepo.GetByEmail(c.Request.Context(), googleUser.Email)
	if err != nil {
		// Create new user
		user = &models.User{
			ID:        utils.GenerateUUID(),
			Email:     googleUser.Email,
			Name:      googleUser.Name,
			Provider:  "google",
			ProviderID: googleUser.ID,
			Picture:   googleUser.Picture,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if err := a.userRepo.Create(c.Request.Context(), user); err != nil {
			utils.Error("Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		// Update user info
		user.Name = googleUser.Name
		user.Picture = googleUser.Picture
		user.UpdatedAt = time.Now()
		
		if err := a.userRepo.Update(c.Request.Context(), user); err != nil {
			utils.Error("Failed to update user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}
	
	// Generate JWT tokens
	accessToken, err := a.generateJWT(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		utils.Error("Failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	refreshToken, err := a.generateJWT(user.ID, user.Email, 7*24*time.Hour)
	if err != nil {
		utils.Error("Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	// Clear OAuth state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	
	c.JSON(http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
		ExpiresIn:    86400, // 24 hours in seconds
	})
}

// RefreshToken refreshes the JWT token
// @Summary Refresh JWT token
// @Description Refresh an expired JWT token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body map[string]string true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (a *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// Verify refresh token
	claims, err := a.verifyJWT(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	
	userID := claims["user_id"].(string)
	email := claims["email"].(string)
	
	// Get user
	user, err := a.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	
	// Generate new tokens
	accessToken, err := a.generateJWT(userID, email, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	newRefreshToken, err := a.generateJWT(userID, email, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	c.JSON(http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		User:         *user,
		ExpiresIn:    86400,
	})
}

// Logout logs out the user
// @Summary Logout user
// @Description Logout the current user
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (a *AuthController) Logout(c *gin.Context) {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetProfile gets the current user profile
// @Summary Get user profile
// @Description Get the current authenticated user's profile
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/profile [get]
func (a *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	user, err := a.userRepo.GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// generateJWT generates a JWT token
func (a *AuthController) generateJWT(userID, email string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	})
	
	return token.SignedString([]byte(a.jwtSecret))
}

// verifyJWT verifies a JWT token
func (a *AuthController) verifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
} 