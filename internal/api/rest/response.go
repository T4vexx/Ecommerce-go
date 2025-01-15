package rest

import "github.com/gofiber/fiber/v2"

func ErrorMessage(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(err.Error())
}

func SuccessMessage(ctx *fiber.Ctx, msg string, data interface{}) error {
	return ctx.Status(200).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}
