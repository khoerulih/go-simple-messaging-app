package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/khoerulih/go-simple-messaging-app/app/repository"
	"github.com/khoerulih/go-simple-messaging-app/pkg/response"
	"go.elastic.co/apm"
)

func GetMessageHistory(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "GetMessageHistory", "controller")
	defer span.End()

	resp, err := repository.GetAllMessages(spanCtx)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada server", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
