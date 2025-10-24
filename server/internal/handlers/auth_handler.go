package handlers

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/database"
	"github.com/Pranavp37/magic_movie_stream/internal/middleware"
	"github.com/Pranavp37/magic_movie_stream/internal/models"
	"github.com/Pranavp37/magic_movie_stream/internal/utils"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	// db := h.handler.Database("magic_movie_stream").Collection("User")
	db := database.GetMongoCollection("magic_movie_stream", "User")

	var userData models.UserRegister

	if err := c.ShouldBindJSON(&userData); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "name , email , passwords are reqried"})
		log.Printf("Error: %v\nStack Trace:\n%s", err, string(debug.Stack()))
		return

	}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	logger.Info("data received ", userData)

	EmailIsExit, err := utils.EmailIsExit(ctx, userData.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error "})
		log.Printf("Error: %v\nStack Trace:\n%s", err, string(debug.Stack()))
		return
	}

	if EmailIsExit {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		log.Printf("Error: %v\nStack Trace:\n%s", err, string(debug.Stack()))
		return
	}
	pass := []byte(userData.Password)
	hashedPassword, err := utils.PasswordHashing(pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		log.Printf("Error: %v\nStack Trace:\n%s", err, string(debug.Stack()))
		return
	}

	Id := utils.GenerateId()

	token, err := middleware.GenerateAccessAndRefreshToken(Id, userData.Email, userData.Name, userData.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		log.Printf("Error: %v\nStack Trace:\n%s", err, string(debug.Stack()))
		return

	}

	user := models.UserRegister{
		User_id:      Id,
		Name:         userData.Name,
		PhomeNumber:  userData.PhomeNumber,
		Email:        userData.Email,
		Role:         userData.Role,
		Password:     hashedPassword,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	_, err = db.InsertOne(ctx, user)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "hard to find collections"})
		log.Printf("Error: %v\nStack Trace:\n%s", err, string(debug.Stack()))
		return

	}

	response := models.UserRegister{
		User_id:      Id,
		Name:         userData.Name,
		PhomeNumber:  user.PhomeNumber,
		Role:         userData.Role,
		Email:        userData.Email,
		Password:     hashedPassword,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "data added successfully",
		"data":    response,
	})

}
