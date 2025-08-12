package repository

import (
	"context"
	"fmt"

	"github.com/khoerulih/go-simple-messaging-app/app/models"
	"github.com/khoerulih/go-simple-messaging-app/pkg/database"
	"go.elastic.co/apm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertNewMessage(ctx context.Context, data models.MessagePayload) error {
	span, _ := apm.StartSpan(ctx, "InsertNewMessage", "repository")
	defer span.End()

	_, err := database.MongoDB.InsertOne(ctx, data)
	return err
}

func GetAllMessages(ctx context.Context) ([]models.MessagePayload, error) {
	span, _ := apm.StartSpan(ctx, "GetAllMessages", "repository")
	defer span.End()

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
