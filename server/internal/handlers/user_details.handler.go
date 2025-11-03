package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/database"
	"github.com/Pranavp37/magic_movie_stream/internal/middleware"
	"github.com/Pranavp37/magic_movie_stream/internal/models"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserDetails(c *gin.Context) {

	token := middleware.JwtTokenFromHeader(c)

	cliem, err := middleware.GetTokenData(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Somthing went wrong"})
		logger.Error("fail to get user details" + err.Error())
	}

	dbCollection := database.GetMongoCollection("chat_application", "user")

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": cliem.User_id,
	}

	user := models.UserRegister{}

	err = dbCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := gin.H{
		"user_id": user.User_id,
		"email":   user.Email,
		"name":    user.Name,
		"phone":   user.PhoneNumber,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
