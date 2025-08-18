package service

import (
	"fmt"
	"encoding/json"
	"practise/go_fiber/internal/database"
	"practise/go_fiber/internal/models"
	"github.com/gofiber/fiber/v2"
)

var SampleArray = make([]interface{}, 0)

func ServerStatus(ctx *fiber.Ctx) error {
	res := models.GetApiResponse("api.server.status", "OK", "Server is Alive..!")
	return ctx.JSON(res)
}

func ListEmployees(ctx *fiber.Ctx) error {
	var employees []models.Employees
	result := database.DB.Find(&employees)
	if result.Error != nil {
		res := models.GetApiResponse("api.server.list.error", "FAILED", result.Error.Error())
		return ctx.Status(500).JSON(res)
	}
	res := models.GetApiResponse("api.server.list.employees", "OK", employees)
	return ctx.JSON(res)
}

func AddEmployee(ctx *fiber.Ctx) error {
	request := new(models.ApiRequest)
	if err := ctx.BodyParser(request); err != nil {
		fmt.Println("error: ", err)
		res := models.GetApiResponse("api.add", "ERROR", ctx.Status(400).JSON(fiber.Map{"error": err.Error()}))
		return ctx.JSON(res)
	}
	var emp models.Employees
	if err := json.Unmarshal(request.Request, &emp); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	result := database.DB.Create(&emp)
	if result.Error != nil {
		res := models.GetApiResponse("api.add", "ERROR", result.Error.Error())
		return ctx.Status(500).JSON(res)
	}
	res := models.GetApiResponse("api.add.employee", "OK", "Inserted Record Successfully")
	return ctx.JSON(res)
}

func GetEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var emp models.Employees
	result := database.DB.First(&emp, id)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	res := models.GetApiResponse("api.get.employee", "OK",emp)
	return ctx.JSON(res)
}

func DeleteEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var emp models.Employees
	result := database.DB.Delete(&emp, id)
	if result.RowsAffected == 0 {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	} 
	if result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Unable to delete user"})
	}
	res := models.GetApiResponse("api.get.employee", "OK","Record Deleted Successfully")
	return ctx.JSON(res)
}


func UpdateEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var emp models.Employees
	result := database.DB.First(&emp, id)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	request := new(models.ApiRequest)
	if err := ctx.BodyParser(request); err != nil {
		fmt.Println("error: ", err)
		res := models.GetApiResponse("api.add", "ERROR", ctx.Status(400).JSON(fiber.Map{"error": err.Error()}))
		return ctx.JSON(res)
	}
	var employee models.Employees
	if err := json.Unmarshal(request.Request, &employee); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	emp.FirstName = employee.FirstName
	emp.LastName = employee.LastName
	emp.Email = employee.Email
	emp.Salary = employee.Salary
	database.DB.Save(&emp)
	res := models.GetApiResponse("api.get.employee", "OK","Record Updated Successfully")
	return ctx.JSON(res)
}