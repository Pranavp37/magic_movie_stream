package handlers

import (
	"github.com/gin-gonic/gin"
)

type FileRequest struct {
	Name string `json:"name" binding:"required"`
	File byte   `json:"file" binding:"required"`
}

func FileHandler(c *gin.Context) {

	name := c.PostForm("name")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dir := "/uploadfile/" + file.Filename

	if err := c.SaveUploadedFile(file, dir); err != nil {
		c.JSON(400, gin.H{"error": "failed to save file"})
	}

	// Process the file upload
	c.JSON(200, gin.H{"message": "File uploaded successfully",
		"name": name,
		"file": file.Filename, "path": dir})
}
