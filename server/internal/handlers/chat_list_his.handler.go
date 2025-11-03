package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/database"
	"github.com/Pranavp37/magic_movie_stream/internal/middleware"
	"github.com/Pranavp37/magic_movie_stream/internal/models"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ChatListHistory(c *gin.Context) {

	token := middleware.JwtTokenFromHeader(c)
	cliem, err := middleware.GetTokenData(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Somthing went wrong"})
		logger.Error("fail to get user details" + err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	chatList := database.GetMongoCollection("chat_application", "conversation")

	filter := bson.M{"participants": bson.M{"$in": []string{cliem.User_id}}}

	ots := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"last_updated": -1})

	cursor, err := chatList.Find(ctx, filter, ots)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch conversations"})
		logger.Error("failed to fetch user details", err.Error())
		return
	}
	var data []models.ChatListModel
	if err := cursor.All(ctx, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode conversations"})
		logger.Error("failed to decode conversations" + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limi":   limit,
		"data":   data,
	})

}
