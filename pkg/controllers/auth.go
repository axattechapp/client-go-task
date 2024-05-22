package controllers

import (
	db "client_task/pkg/common/db/sqlc"
	"client_task/pkg/payloads"
	"net/http"

	"client_task/pkg/utils/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (cc *UsersController) LoginHandler(ctx *gin.Context) {

	var loginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.BindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := cc.DB.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		// Handle error, e.g., user not found
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		// Handle invalid password
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := token.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "status": "Login successfully"})
}

func (cc *UsersController) RegisterHandler(c *gin.Context) {
	var payload *payloads.CreateUser

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(payload.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.CreateUserParams{
		FullName:     payload.FullName,
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		UserType:     payload.UserType,
	}

	user, err := cc.DB.CreateUser(c, *args)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "Registration failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Registration successful", "user": user})

}
