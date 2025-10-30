package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/database"
	"github.com/Pranavp37/magic_movie_stream/internal/models"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SearchUser(c *gin.Context) {

	db := database.GetMongoCollection("chat_application", "user")

	query := c.Query("query")

	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query 'param query' is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	filter := bson.M{
		"name": bson.M{
			"$regex":   query,
			"$options": "i", // ✅ notice the "s" — correct MongoDB operator
		},
	}

	cusor, err := db.Find(ctx, filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	defer cusor.Close(ctx)

	var users []models.UserSearchResponse
	if err := cusor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"response": users,
	})

}
