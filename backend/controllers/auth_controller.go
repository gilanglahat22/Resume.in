package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	oauthConfig       *oauth2.Config
	oauthRegConfig    *oauth2.Config
	userRepo          models.UserRepository
	jwtSecret         string
}

// RegistrationRequest represents the registration request body
type RegistrationRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
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
		oauthRegConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/user.birthday.read",
				"https://www.googleapis.com/auth/user.gender.read",
				"https://www.googleapis.com/auth/user.phonenumbers.read",
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
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Birthday      string `json:"birthday,omitempty"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
	ExpiresIn    int64       `json:"expires_in"`
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
// @Description Processes the OAuth callback and returns JWT tokens (handles both login and registration)
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
	
	// Check if this is a registration flow
	isRegistration := strings.HasPrefix(state, "register_")
	
	// Verify state based on flow type
	var storedState string
	var err error
	if isRegistration {
		storedState, err = c.Cookie("oauth_reg_state")
	} else {
		storedState, err = c.Cookie("oauth_state")
	}
	
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	
	// Choose the appropriate config based on flow type
	var oauthConfig *oauth2.Config
	if isRegistration {
		// Use registration config for enhanced scopes
		oauthConfig = a.oauthRegConfig
	} else {
		oauthConfig = a.oauthConfig
	}
	
	// Exchange code for token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.Error("Failed to exchange token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	
	// Get user info from Google
	client := oauthConfig.Client(context.Background(), token)
	
	// Get basic profile info
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
	
	// Get additional user information if registration flow
	if isRegistration {
		peopleResp, err := client.Get("https://people.googleapis.com/v1/people/me?personFields=phoneNumbers,birthdays,genders")
		if err == nil {
			defer peopleResp.Body.Close()
			var peopleInfo struct {
				PhoneNumbers []struct {
					Value string `json:"value"`
				} `json:"phoneNumbers"`
				Birthdays []struct {
					Date struct {
						Year  int `json:"year"`
						Month int `json:"month"`
						Day   int `json:"day"`
					} `json:"date"`
				} `json:"birthdays"`
				Genders []struct {
					Value string `json:"value"`
				} `json:"genders"`
			}
			
			if err := json.NewDecoder(peopleResp.Body).Decode(&peopleInfo); err == nil {
				if len(peopleInfo.PhoneNumbers) > 0 {
					googleUser.PhoneNumber = peopleInfo.PhoneNumbers[0].Value
				}
				if len(peopleInfo.Genders) > 0 {
					googleUser.Gender = peopleInfo.Genders[0].Value
				}
				if len(peopleInfo.Birthdays) > 0 {
					date := peopleInfo.Birthdays[0].Date
					googleUser.Birthday = fmt.Sprintf("%04d-%02d-%02d", date.Year, date.Month, date.Day)
				}
			}
		}
	}
	
	// Check if user exists
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
		// If registration flow and user exists, return error
		if isRegistration {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		
		// Update user info for login flow
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
	
	// Clear appropriate OAuth state cookie
	if isRegistration {
		c.SetCookie("oauth_reg_state", "", -1, "/", "", false, true)
	} else {
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	}
	
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

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param registration body RegistrationRequest true "Registration details"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 409 {object} map[string]string "Email already registered"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Check if user already exists
	existingUser, err := a.userRepo.GetByEmail(c.Request.Context(), req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	// Create new user
	user := &models.User{
		ID:        utils.GenerateUUID(),
		Email:     req.Email,
		Name:      req.Name,
		Password:  hashedPassword,
		Provider:  "local",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.userRepo.Create(c.Request.Context(), user); err != nil {
		utils.Error("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
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

	c.JSON(http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
		ExpiresIn:    86400, // 24 hours in seconds
	})
}

// GoogleRegister initiates the Google OAuth registration flow
// @Summary Initiate Google OAuth registration
// @Description Redirects to Google OAuth consent page with additional scopes for registration
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "Returns auth_url for consent screen"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/google/register [get]
func (a *AuthController) GoogleRegister(c *gin.Context) {
	state := "register_" + utils.GenerateRandomString(32)
	
	// Store state in session or temporary storage for verification
	c.SetCookie("oauth_reg_state", state, 600, "/", "", false, true)
	
	// Use the same redirect URL as login but with different state
	tempConfig := *a.oauthRegConfig
	tempConfig.RedirectURL = a.oauthConfig.RedirectURL // Use same redirect URL as login
	
	url := tempConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
	})
}

// GoogleRegisterCallback handles the Google OAuth registration callback
// @Summary Handle Google OAuth registration callback
// @Description Processes the OAuth callback for registration with additional user information
// @Tags auth
// @Param code query string true "Authorization code"
// @Param state query string true "OAuth state"
// @Produce json
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/google/register/callback [get]
func (a *AuthController) GoogleRegisterCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	
	// Verify state
	storedState, err := c.Cookie("oauth_reg_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	
	// Exchange code for token
	token, err := a.oauthRegConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.Error("Failed to exchange token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	
	// Get user info from Google
	client := a.oauthRegConfig.Client(context.Background(), token)
	
	// Get basic profile info
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
	
	// Get additional user information
	peopleResp, err := client.Get("https://people.googleapis.com/v1/people/me?personFields=phoneNumbers,birthdays,genders")
	if err == nil {
		defer peopleResp.Body.Close()
		var peopleInfo struct {
			PhoneNumbers []struct {
				Value string `json:"value"`
			} `json:"phoneNumbers"`
			Birthdays []struct {
				Date struct {
					Year  int `json:"year"`
					Month int `json:"month"`
					Day   int `json:"day"`
				} `json:"date"`
			} `json:"birthdays"`
			Genders []struct {
				Value string `json:"value"`
			} `json:"genders"`
		}
		
		if err := json.NewDecoder(peopleResp.Body).Decode(&peopleInfo); err == nil {
			if len(peopleInfo.PhoneNumbers) > 0 {
				googleUser.PhoneNumber = peopleInfo.PhoneNumbers[0].Value
			}
			if len(peopleInfo.Genders) > 0 {
				googleUser.Gender = peopleInfo.Genders[0].Value
			}
			if len(peopleInfo.Birthdays) > 0 {
				date := peopleInfo.Birthdays[0].Date
				googleUser.Birthday = fmt.Sprintf("%04d-%02d-%02d", date.Year, date.Month, date.Day)
			}
		}
	}
	
	// Check if user exists
	existingUser, err := a.userRepo.GetByEmail(c.Request.Context(), googleUser.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	
	// Create new user
	user := &models.User{
		ID:         utils.GenerateUUID(),
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		Provider:   "google",
		ProviderID: googleUser.ID,
		Picture:    googleUser.Picture,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := a.userRepo.Create(c.Request.Context(), user); err != nil {
		utils.Error("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
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
	c.SetCookie("oauth_reg_state", "", -1, "/", "", false, true)
	
	c.JSON(http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
		ExpiresIn:    86400, // 24 hours in seconds
	})
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