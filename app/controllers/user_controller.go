package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khoerulih/go-simple-messaging-app/app/models"
	"github.com/khoerulih/go-simple-messaging-app/app/repository"
	jwttoken "github.com/khoerulih/go-simple-messaging-app/pkg/jwt_token"
	"github.com/khoerulih/go-simple-messaging-app/pkg/response"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		errResponse := fmt.Errorf("failed to parse request body: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	if err := user.Validate(); err != nil {
		errResponse := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to hash password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	user.Password = string(hashedPassword)

	if err := repository.InsertNewUser(spanCtx, user); err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)

	}

	resp := user
	resp.Password = ""

	return response.SendSuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	// parsing and validating request
	loginReq := new(models.LoginRequest)
	loginResp := models.LoginResponse{}
	now := time.Now()

	if err := ctx.BodyParser(loginReq); err != nil {
		errResponse := fmt.Errorf("failed to parse request body: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	if err := loginReq.Validate(); err != nil {
		errResponse := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	user, err := repository.GetUserByUsername(spanCtx, loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get user by username: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username atau password salah", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		errResponse := fmt.Errorf("failed to check password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username atau password salah", nil)
	}

	token, err := jwttoken.GenerateToken(spanCtx, user.Username, user.Fullname, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	refreshToken, err := jwttoken.GenerateToken(spanCtx, user.Username, user.Fullname, "refresh_token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate refresh token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	userSession := &models.UserSession{
		UserID:              int(user.ID),
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(jwttoken.MapTokenType["token"]),
		RefreshTokenExpired: now.Add(jwttoken.MapTokenType["refresh_token"]),
	}

	err = repository.InsertNewUserSession(spanCtx, userSession)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user session: %w", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	loginResp.Username = user.Username
	loginResp.Fullname = user.Fullname
	loginResp.Token = token
	loginResp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, loginResp)
}

func Logout(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Logout", "controller")
	defer span.End()

	token := ctx.Get("Authorization")

	err := repository.DeleteUserSessionByToken(spanCtx, token)
	if err != nil {
		errResponse := fmt.Errorf("failed to delete user session: %w", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}
	return response.SendSuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	now := time.Now()
	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwttoken.GenerateToken(spanCtx, username, fullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %w", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	err = repository.UpdateUserSession(spanCtx, token, now.Add(jwttoken.MapTokenType["token"]), refreshToken)
	if err != nil {
		errResponse := fmt.Errorf("failed to updatetoken: %w", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}
	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
