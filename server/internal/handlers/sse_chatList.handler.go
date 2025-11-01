package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Pranavp37/magic_movie_stream/internal/models"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

var sseClients = make(map[string]chan string)
var sseMux sync.Mutex

func ChatListSSE(c *gin.Context) {

	user_id := c.Query("userID")

	if user_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userID is requied",
		})
		logger.Error("userID is required")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	msgChan := make(chan string)

	sseMux.Lock()
	sseClients[user_id] = msgChan
	sseMux.Unlock()
	defer func() {
		sseMux.Lock()
		delete(sseClients, user_id)
		sseMux.Unlock()
		close(msgChan)
	}()

	// Send initial connection message
	fmt.Fprintf(c.Writer, "data: {\"type\":\"connected\"}\n\n")
	c.Writer.Flush()

	clientGone := c.Request.Context().Done()

	for {
		select {
		case <-clientGone:
			// Client disconnected
			logger.Infof("SSE client disconnected: %s", user_id)
			return
		case msg, ok := <-msgChan:
			if !ok {
				return
			}

			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			c.Writer.Flush()

		}
	}

}

func NotifySSEChatListUpdate(msg *models.MessageReceivedModel, conversationID string, receverID string, senderID string) {

	updatedData := map[string]interface{}{
		"conversation_id": conversationID,
		"participants":    []string{senderID, receverID},
		"last_message":    msg.Message,
		"last_updated":    msg.TimeStamp,
	}

	jsonData, _ := json.Marshal(updatedData)

	for _, userID := range []string{senderID, receverID} {
		sseMux.Lock()
		ch, ok := sseClients[userID]
		sseMux.Unlock()

		if ok {
			select {
			case ch <- string(jsonData):
				logger.Infof("SSE notification sent to user %s", userID)
			default:
				logger.Warnf("SSE channel full for user %s", userID)
			}

		}
	}

}
