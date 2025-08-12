package repository

import (
	"context"
	"fmt"

	"github.com/khoerulih/go-simple-messaging-app/app/models"
	"github.com/khoerulih/go-simple-messaging-app/pkg/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertMessage(ctx context.Context, data models.MessagePayload) error {
	_, err := database.MongoDB.InsertOne(ctx, data)
	return fmt.Errorf("failed to insert new message: %v", err)
}

func GetAllMessages(ctx context.Context) ([]models.MessagePayload, error) {
	var (
		err  error
		resp []models.MessagePayload
	)

	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %v", err)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var payload models.MessagePayload
		if err := cursor.Decode(&payload); err != nil {
			return nil, fmt.Errorf("failed to decode message: %v", err)
		}
		resp = append(resp, payload)
	}

	return resp, nil
}
