package service

import (
	"practise/go_fiber/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ServerStatus(ctx *fiber.Ctx) error {
	res := models.GetApiResponse("api.server.status", "OK", "Server is Alive..!")
	return ctx.JSON(res)
}

