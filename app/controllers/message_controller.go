package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/khoerulih/go-simple-messaging-app/app/repository"
	"github.com/khoerulih/go-simple-messaging-app/pkg/response"
)

func GetMessageHistory(ctx *fiber.Ctx) error {
	resp, err := repository.GetAllMessages(ctx.Context())
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada server", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
