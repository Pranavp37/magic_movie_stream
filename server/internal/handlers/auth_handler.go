package handlers

import (
	"context"
	"fmt"
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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(c *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("internal error: %v", r),
			})
		}
	}()
	// db := h.handler.Database("magic_movie_stream").Collection("User")
	db := database.GetMongoCollection("chat_application", "user")
	if db == nil {
		logger.Error("Mongo collection is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

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

	user := models.UserRegisterResponse{
		User_id:      Id,
		Name:         userData.Name,
		PhoneNumber:  userData.PhoneNumber,
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

	response := models.UserRegisterResponse{
		User_id:      Id,
		Name:         userData.Name,
		PhoneNumber:  user.PhoneNumber,
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

func Login(c *gin.Context) {
	db := database.GetMongoCollection("chat_application", "user")

	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email and password are required",
		})
		logger.Error("email and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	filter := bson.M{
		"email": req.Email,
	}

	var user models.UserRegister
	err := db.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		logger.Error("user not found")
		return

	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Error(err.Error())
		return
	}

	err = utils.PasswordCompare([]byte(req.Password), []byte(user.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		logger.Error("invalid password")
		return
	}

	token, err := middleware.GenerateAccessAndRefreshToken(user.User_id, user.Email, user.Name, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't generate token",
		})
		logger.Error("can't generate token")
		return
	}

	response := models.UserLoginResponse{
		User_id:      user.User_id,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Role:         user.Role,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})

}
