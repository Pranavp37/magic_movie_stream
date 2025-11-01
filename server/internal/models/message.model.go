package models

import (
	"time"
)

// MessageType defines the type of message (e.g., text, image, video)
type MessageType string

const (
	TextMessage   MessageType = "text"
	ImageMessage  MessageType = "image"
	VideoMessage  MessageType = "video"
	AudioMessage  MessageType = "audio"
	FileMessage   MessageType = "file"
	SystemMessage MessageType = "system"
	MessageSeen   MessageType = "message_seen"
	UserStatus    MessageType = "user_status"
)

type MessageReceivedModel struct {
	ConversationID string      `json:"conversation_id" bson:"conversation_id"`
	MessageID      string      `json:"message_id" bson:"message_id"`
	SenderID       string      `json:"sender_id" bson:"sender_id"`
	ReceiverID     string      `json:"receiver_id" bson:"receiver_id"`
	Type           MessageType `json:"message_type" bson:"message_type"`
	Message        string      `json:"message" bson:"message"`
	TimeStamp      time.Time   `json:"timestamp" bson:"timestamp"`
	Seen           bool        `json:"seen" bson:"seen"`
}

type ChatListModel struct {
	ConversationID string    `json:"conversation_id" bson:"conversation_id"`
	Participants   []string  `json:"participants" bson:"participants"`
	LastMessage    string    `json:"last_message" bson:"last_message"`
	LastUpdated    time.Time `json:"last_updated" bson:"last_updated"`
}

type UserStatuSendToUser struct {
	Type   MessageType `json:"message_type" bson:"message_type"`
	UserId string      `json:"user_id" bson:"user_id"`
	Online bool        `json:"isOnline" bson:"isOnline"`
}
