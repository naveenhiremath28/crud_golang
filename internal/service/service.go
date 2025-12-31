package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"practise/go_fiber/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) ServerStatus(ctx *fiber.Ctx) error {
	res := models.GetApiResponse("api.server.status", "OK", "Server is Alive..!")
	return ctx.JSON(res)
}

func (s *Service) ListEmployees(ctx *fiber.Ctx) error {
	var employees []models.Employees
	result := s.DB.Find(&employees)
	if result.Error != nil {
		res := models.GetApiResponse("api.server.list.error", "FAILED", result.Error.Error())
		return ctx.Status(500).JSON(res)
	}
	res := models.GetApiResponse("api.server.list.employees", "OK", employees)
	return ctx.JSON(res)
}

func (s *Service) AddEmployee(ctx *fiber.Ctx) error {
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
	if emp.ID == "" {
		emp.ID = uuid.New().String()
	}
	emp.VaultEntityID = emp.ID + "pii"
	result := s.DB.Create(&emp)
	if result.Error != nil {
		res := models.GetApiResponse("api.add", "ERROR", result.Error.Error())
		return ctx.Status(500).JSON(res)
	}
	res := models.GetApiResponse("api.add.employee", "OK", "Inserted Record Successfully")
	return ctx.JSON(res)
}

func (s *Service) GetEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var emp models.Employees
	result := s.DB.First(&emp, "id = ?", id)
	if result.Error != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	res := models.GetApiResponse("api.get.employee", "OK", emp)
	return ctx.JSON(res)
}

func (s *Service) DeleteEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var emp models.Employees
	result := s.DB.Delete(&emp, "id = ?", id)
	if result.RowsAffected == 0 {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	if result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Unable to delete user"})
	}
	res := models.GetApiResponse("api.get.employee", "OK", "Record Deleted Successfully")
	return ctx.JSON(res)
}

func (s *Service) UpdateEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var emp models.Employees
	result := s.DB.First(&emp, "id = ?", id)
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
	emp.VaultEntityID = emp.ID + "pii"
	s.DB.Save(&emp)
	res := models.GetApiResponse("api.get.employee", "OK", "Record Updated Successfully")
	return ctx.JSON(res)
}

func (s *Service) LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest[models.Login]
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	form := url.Values{}
	form.Set("client_id", "employee-api")
	form.Set("grant_type", "password")
	form.Set("username", req.Request.Username)
	form.Set("password", req.Request.Password)

	resp, err := http.PostForm(
		"http://localhost:8083/realms/employee_realm/protocol/openid-connect/token",
		form,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to reach keycloak",
		})
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse keycloak response",
		})
	}

	// You can wrap this with your GetApiResponse if you want later
	return c.Status(resp.StatusCode).JSON(data)
}

func (s *Service) RefreshHandler(c *fiber.Ctx) error {
	var req models.LoginRequest[models.Refresh]
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Request.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh_token is required",
		})
	}

	form := url.Values{}
	form.Set("client_id", "employee-api")
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", req.Request.RefreshToken)

	resp, err := http.PostForm(
		"http://localhost:8083/realms/employee_realm/protocol/openid-connect/token",
		form,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to reach keycloak",
		})
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse keycloak response",
		})
	}

	return c.Status(resp.StatusCode).JSON(data)
}
