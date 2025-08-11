package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	SuccessMessage = "success"
)

func SendSuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(Response{
		Message: SuccessMessage,
		Data:    data,
	})
}

func SendFailureResponse(ctx *fiber.Ctx, httpcode int, message string, data interface{}) error {
	return ctx.Status(httpcode).JSON(Response{
		Message: message,
		Data:    data,
	})
}
