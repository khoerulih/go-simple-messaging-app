package router

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khoerulih/go-simple-messaging-app/app/repository"
	jwttoken "github.com/khoerulih/go-simple-messaging-app/pkg/jwt_token"
	"github.com/khoerulih/go-simple-messaging-app/pkg/response"
)

func MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		fmt.Println("authorization empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	_, err := repository.GetUserSessionByToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println("failed to get user session on DB: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwttoken.ValidateToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		fmt.Println("jwt token is expired : ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		fmt.Println("authorization empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwttoken.ValidateToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		fmt.Println("jwt token is expired : ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}
