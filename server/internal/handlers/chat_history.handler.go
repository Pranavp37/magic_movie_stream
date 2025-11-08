package handlers

import (
	"context"
	"fmt"
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

func ChatHistory(c *gin.Context) {

	conv_ID := c.Query("conversationID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	skip := (page - 1) * limit

	token := middleware.JwtTokenFromHeader(c)

	_, err := middleware.GetTokenData(token)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Somthing went wrong..."})
		logger.Error("fail to get user details" + err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filters := bson.M{
		"conversation_id": conv_ID,
	}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(bson.M{"created_at": -1})
	chatDB := database.GetMongoCollection("chat_application", "message")
	cursor, err := chatDB.Find(ctx, filters, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch chat history"})
		logger.Error("failed to fetch chat history" + err.Error())
		return
	}
	var data []models.MessageReceivedModel
	if err := cursor.All(ctx, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "somthing went wrong..."})
		logger.Error("failed to fetch chat history")
		return
	}

	fmt.Printf("KKKKK>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%v\n", data)
	c.JSON(http.StatusOK, gin.H{"status": "success",
		"page":  page,
		"limit": limit,
		"data":  data})
}
