package utils

import (
	"github.com/google/uuid"
)

func GenerateId() string {
	return uuid.New().String()
}

func GenerateConversationId(senderID, receiverID string) string {

	if senderID == "" || receiverID == "" {
		return ""
	}

	return "conv_" + senderID + "_" + receiverID

}
