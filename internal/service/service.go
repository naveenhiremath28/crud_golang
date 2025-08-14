package service

import (
	"fmt"
	"practise/go_fiber/internal/database"
	"practise/go_fiber/internal/models"
	"github.com/gofiber/fiber/v2"
)

var SampleArray = make([]interface{}, 0)

func ServerStatus(ctx *fiber.Ctx) error {
	res := models.GetApiResponse("api.server.status", "OK", "Server is Alive..!")
	return ctx.JSON(res)
}

func List(ctx *fiber.Ctx) error {
	// var greeting string
	// err := database.DB.QueryRow("SELECT * from employees;").Scan(&greeting)
	// if err != nil {
	// 	return err
	// }
	var employees []models.Employees
	database.DB.Find(&employees)
	res := models.GetApiResponse("api.server.list.array", "OK", employees)
	return ctx.JSON(res)
}


func Add(ctx *fiber.Ctx) error {
	request := new(models.ApiRequest)
	if err := ctx.BodyParser(request); err != nil {
		fmt.Println("error: ", err)
		res := models.GetApiResponse("api.add", "ERROR", ctx.Status(400).JSON(fiber.Map{"error": err.Error()}))
		return ctx.JSON(res)
	}
	fmt.Println("request: ", request.Request)
	// database.DB.Create(&employees)
	res := models.GetApiResponse("api.server.status", "OK", "Not IMPLEMENTED")
	return ctx.JSON(res)
}