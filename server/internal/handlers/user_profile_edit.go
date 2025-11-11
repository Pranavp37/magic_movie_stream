package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/database"
	"github.com/Pranavp37/magic_movie_stream/internal/utils"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UserProfileEdit(c *gin.Context) {

	name := c.PostForm("name")
	phone := c.PostForm("phone_number")
	email := c.PostForm("email")
	profilePic, _ := c.FormFile("profile_pic")
	id := c.PostForm("user_id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		logger.Error("user_id is required")
		return

	}
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	var profilePicURL string
	if profilePic != nil {
		Uploader, err := utils.NewS3Uploader(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize S3 uploader"})
			logger.Error("S3 UPLOADER error" + err.Error())
			return
		}

		url, err := Uploader.FileUploader(ctx, profilePic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "file upload failed"})
			logger.Error("file upload failed" + err.Error())
			return
		}

		profilePicURL = url
	} else {
		profilePicURL = ""
	}

	db := database.GetMongoCollection("chat_application", "user")
	if db == nil {
		logger.Error("Mongo collection is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	filters := bson.M{
		"user_id": id,
	}

	updateFields := bson.M{}
	if name != "" {
		updateFields["name"] = name
	}
	if phone != "" {
		updateFields["phone_number"] = phone
	}
	if email != "" {
		updateFields["email"] = email
	}
	if profilePicURL != "" {
		updateFields["profile_pic"] = profilePicURL
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields to update"})
		return
	}
	update := bson.M{"$set": updateFields}
	_, err := db.UpdateOne(ctx, filters, update)
	if err != nil {
		logger.Error("Failed to update user profile: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
