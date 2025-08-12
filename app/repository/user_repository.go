package repository

import (
	"context"
	"time"

	"github.com/khoerulih/go-simple-messaging-app/app/models"
	"github.com/khoerulih/go-simple-messaging-app/pkg/database"
	"go.elastic.co/apm"
)

func InsertNewUser(ctx context.Context, user *models.User) error {
	span, _ := apm.StartSpan(ctx, "InsertNewUser", "repository")
	defer span.End()
	return database.DB.Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	span, _ := apm.StartSpan(ctx, "GetUserByUsername", "repository")
	defer span.End()
	var (
		resp models.User
		err  error
	)
	err = database.DB.Where("username = ?", username).Last(&resp).Error
	return resp, err
}

func InsertNewUserSession(ctx context.Context, session *models.UserSession) error {
	span, _ := apm.StartSpan(ctx, "InsertNewUserSession", "repository")
	defer span.End()
	return database.DB.Create(session).Error
}

func GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error) {
	span, _ := apm.StartSpan(ctx, "GetUserSessionByToken", "repository")
	defer span.End()

	var (
		resp models.UserSession
		err  error
	)
	err = database.DB.Where("token = ?", token).Last(&resp).Error
	return resp, err
}

func UpdateUserSession(ctx context.Context, token string, tokenExpired time.Time, refreshToken string) error {
	span, _ := apm.StartSpan(ctx, "UpdateUserSession", "repository")
	defer span.End()
	return database.DB.Exec(`UPDATE user_sessions SET token = ?, token_expired = ? WHERE refresh_token = ?`, token, tokenExpired, refreshToken).Error
}

func DeleteUserSessionByToken(ctx context.Context, token string) error {
	span, _ := apm.StartSpan(ctx, "DeleteUserSessionByToken", "repository")
	defer span.End()
	return database.DB.Exec(`DELETE FROM user_sessions WHERE token = ?`, token).Error
}
