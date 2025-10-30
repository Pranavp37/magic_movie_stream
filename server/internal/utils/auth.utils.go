package utils

import (
	"context"

	"github.com/Pranavp37/magic_movie_stream/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func EmailIsExit(ctx context.Context, email string) (bool, error) {

	collection := database.GetMongoCollection("chat_application", "user")

	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return false, nil
		}
		return false, err
	}
	return true, nil
}

func PasswordHashing(password []byte) (string, error) {

	hashedpassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedpassword), nil
}

func PasswordCompare(password, hashedPassword []byte) error {

	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return err

	}

	return nil
}
