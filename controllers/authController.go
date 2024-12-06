package controllers

import (
	"log"
	"net/http"
	"tender_management/models"
	"tender_management/pkg/db/password"
	"tender_management/pkg/redise"
	"tender_management/pkg/utils"
	"tender_management/validation"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthController struct {
	Storage *gorm.DB
	Redis   *redise.RedisDB
}

func NewAuthController(storage *gorm.DB, redis *redise.RedisDB) *AuthController {
	return &AuthController{
		Storage: storage,
		Redis:   redis,
	}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Registers a new user by providing phone number, email, and password.
//
//	Validates phone number, email, and password, then generates a verification code for phone number validation.
//
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.UserRegister true "User Registration Data"
// @Success      201 {object} Response "Successfully created the user"
// @Failure      400 {object} Response "Bad request"
// @Failure      500 {object} Response "Internal server error"
// @Router       /auth/register [post]
func (ac *AuthController) CreateUser(c *gin.Context) {
	var user models.UserRegister

	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request user", err)
		return
	}

	if !validation.IsValidPhoneNumber(user.PhoneNumber) {
		handleError(c, http.StatusBadRequest, "Invalid phone number format. Must start with +998 and be 12 digits long", nil)
		return
	}

	if err := validation.ValidatePassword(user.Password); err != nil {
		handleError(c, http.StatusBadRequest, "Password validation failed", err)
		return
	}

	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Password hashing failed", err)
		return
	}

	code := utils.GenerateCode(6)

	info := map[string]interface{}{
		"first_name":   user.FirstName,
		"phone_number": user.PhoneNumber,
		"email":        user.Email,
		"password":     hashedPassword,
		"user_role":    user.Role,
		"is_active":    true,
		"code":         code,
	}

	infoJson, err := json.Marshal(info)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to marshal userinfo to JSON", err)
		return
	}

	if err := ac.Redis.SetEx(c, user.PhoneNumber, infoJson, 3*time.Minute); err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to save verification code and user info in Redis", err)
		return
	}

	response := map[string]interface{}{
		"message":      "Verification code sent successfully",
		"phone_number": user.PhoneNumber,
		"expires_in":   3 * time.Minute,
	}

	HandleResponse(c, http.StatusCreated, response)
}

// VerifyCode godoc
// @Summary      Verify phone number with the code
// @Description  Verifies the user's phone number using the code sent earlier. If valid, activates the user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.VerifyRequest true "Phone number and code verification data"
// @Success      200 {string} string "User verified and activated successfully"
// @Failure      400 {object} map[string]interface{} "Invalid or expired code"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /auth/verify [post]
func (ac *AuthController) VerifyCode(c *gin.Context) {
	var body models.VerifyRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request body", err)
		return
	}

	info, err := ac.Redis.Get(c, body.PhoneNumber)
	if err == redis.Nil {
		handleError(c, http.StatusBadRequest, "Verification code not found or expired", err)
		return
	}else if err != nil {
		handleError(c, http.StatusInternalServerError, "Redis server error", err)
		return
	}

	var userinfo map[string]interface{}
	if err := json.Unmarshal([]byte(info), &userinfo); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON format: Failed to parse user information", err)
		return
	}

	if !validation.VerifyCode(userinfo["code"].(string), body.Code) {
		handleError(c, http.StatusBadRequest, "Wrong OTP code, please try again", nil)
		return
	}

	users := models.Users{
		FirstName:   userinfo["first_name"].(string),
		PhoneNumber: userinfo["phone_number"].(string),
		Email:       userinfo["email"].(string),
		Password:    userinfo["password"].(string),
		Role:        userinfo["user_role"].(string),
		IsActive:    true,
	}

	if err := ac.Storage.Create(&users).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	if err := ac.Redis.Delete(c, body.PhoneNumber); err != nil {
		log.Println("Failed to delete Redis key:", body.PhoneNumber)
	}

	HandleResponse(c, http.StatusOK, "User verified and activated successfully")
}

// LoginUser godoc
// @Summary      Login a user
// @Description  Allows a user to log in using email and password. If valid, returns a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginRequest true "Login Credentials"
// @Success      200 {object} Response "Successfully restored the category"
// @Failure      400 {object} Response "Invalid request"
// @Failure      401 {object} Response "Unauthorized: Invalid credentials"
// @Failure      404 {object} Response "User not found"
// @Failure      500 {object} Response "Internal server error"
// @Router       /auth/login [post]
func (ac *AuthController) LoginUser(c *gin.Context) {
	var login models.LoginRequest
	var user models.Users

	if err := c.ShouldBindJSON(&login); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request login", err)
		return
	}

	if err := ac.Storage.Where("email = ?", login.Email).First(&user).Error; err != nil {
		handleError(c, http.StatusNotFound, "User not found", err)
		return
	}

	if !user.IsActive {
		handleError(c, http.StatusUnauthorized, "User is not verified yet", nil)
		return
	}

	if !password.CheckPasswordHash(login.Password, user.Password) {
		handleError(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Role)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	HandleResponse(c, http.StatusOK, gin.H{
		"token": token,
	})
}


