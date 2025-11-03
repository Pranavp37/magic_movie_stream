package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/database"

	"github.com/Pranavp37/magic_movie_stream/internal/models"
	"github.com/Pranavp37/magic_movie_stream/internal/utils"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var webConnection = make(map[string]*websocket.Conn)
var userOnline = make(map[string]bool)
var webmu = &sync.RWMutex{}

func WebSocketConnection(c *gin.Context) {

	user_id := c.Query("userID")
	red := color.New(color.FgRed).SprintFunc()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		logger.Error("error upgrading websocket" + err.Error())
	}
	webmu.Lock()
	webConnection[user_id] = conn
	userOnline[user_id] = true
	logger.Info("user is connected--------------------------------------")
	webmu.Unlock()
	defer func() {
		conn.Close()
		webmu.Lock()
		delete(webConnection, user_id)
		userOnline[user_id] = false
		BroadcastUserStatus(user_id, false)
		webmu.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("error reading message" + err.Error())
			delete(webConnection, user_id)
			break
		}

		fmt.Println(red("message Recived" + string(message)))

		var msg models.MessageReceivedModel
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Error("Invalid JSON message: " + err.Error())
			continue

		}

		// infom user is online expect of send user
		BroadcastUserStatus(user_id, true)

		//send data into database
		go AddMessagesTOCollections(&msg)

		switch msg.Type {
		case models.TextMessage:
			//send data to reciver connection if user is online
			SendToSender(msg.ReceiverID, &msg)
			//send data to sender connection if user is online
			SendTOUserSender(msg.SenderID, &msg)

		case models.MessageSeen:
			//upadate user message collection
			UpdateUserMessageSeen(msg.SenderID, &msg)

		}

	}
}

func SendToSender(user_id string, msg *models.MessageReceivedModel) {
	webmu.RLock()
	conn, ok := webConnection[user_id]
	webmu.RUnlock()

	if ok {
		msgByte, err := json.Marshal(msg)
		if err != nil {
			logger.Error("Failed to marshal message:" + err.Error())
		} else {
			if err := conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
				logger.Error("Error writing message to receiver: " + err.Error())
				webmu.Lock()
				delete(webConnection, user_id)
				userOnline[user_id] = false
				webmu.Unlock()
				BroadcastUserStatus(user_id, false)
			}
		}
	}

}
func SendTOUserSender(user_id string, msg *models.MessageReceivedModel) {
	webmu.RLock()
	conn, ok := webConnection[user_id]
	webmu.RUnlock()

	if ok {
		msgByte, err := json.Marshal(msg)
		if err != nil {
			logger.Error("Failed to marshal message:" + err.Error())
		} else {
			if err := conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
				logger.Error("Error writing message to sender: " + err.Error())

				webmu.Lock()
				delete(webConnection, user_id)
				userOnline[user_id] = false
				webmu.Unlock()
				BroadcastUserStatus(user_id, false)

			}
		}
	}

}

func BroadcastUserStatus(user_id string, online bool) {
	webmu.RLock()
	for uid, conn := range webConnection {
		if uid == user_id {
			continue // don't send status to the same user
		}
		data := &models.UserStatuSendToUser{
			Type:   "user_status",
			UserId: user_id,
			Online: online,
		}
		dataByte, _ := json.Marshal(data)
		conn.WriteMessage(websocket.TextMessage, dataByte)
	}
	webmu.RUnlock()

}

func AddMessagesTOCollections(msg *models.MessageReceivedModel) {

	conver_db := database.GetMongoCollection("chat_application", "conversation")
	message_db := database.GetMongoCollection("chat_application", "message")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	convID := utils.GenerateConversationId(msg.SenderID, msg.ReceiverID)

	//for skipping user sending message id
	msg.MessageID = primitive.NewObjectID().String()
	msg.ConversationID = convID

	_, err := message_db.InsertOne(ctx, msg)

	if err != nil {
		logger.Error("Failed to save message: " + err.Error())
	}

	filters := bson.M{
		"conversation_id": convID,
	}

	update := bson.M{
		"$set": bson.M{
			"participants": []string{msg.SenderID, msg.ReceiverID},
			"last_message": msg.Message,
			"last_updated": msg.TimeStamp,
		},
	}
	_, err = conver_db.UpdateOne(context.Background(), filters, update, options.Update().SetUpsert(true))

	if err != nil {
		logger.Error("Failed to update conversation: " + err.Error())
	}

	go NotifySSEChatListUpdate(msg, convID, msg.ReceiverID, msg.SenderID)

}

func UpdateUserMessageSeen(user_id string, msg *models.MessageReceivedModel) {

	message_db := database.GetMongoCollection("chat_application", "message")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	filters := bson.M{
		"message_id": msg.MessageID,
	}

	update := bson.M{
		"$set": bson.M{
			"seen":    true,
			"seen_by": []string{user_id},
			"seen_at": time.Now(),
		},
	}
	_, err := message_db.UpdateOne(ctx, filters, update)
	if err != nil {
		logger.Error("Failed to update conversation: " + err.Error())

	}
}
