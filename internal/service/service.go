package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	DB     *gorm.DB
	Config *config.Config
	Log    *zap.SugaredLogger
}

func NewService(db *gorm.DB, cfg *config.Config, log *zap.SugaredLogger) *Service {
	return &Service{DB: db, Config: cfg, Log: log}
}

func (s *Service) ServerStatus(ctx *fiber.Ctx) error {
	s.Log.Debug("Server status check")
	res := models.GetApiResponse("api.server.status", "OK", "Server is Alive..!")
	return ctx.JSON(res)
}

func (s *Service) ListEmployees(ctx *fiber.Ctx) error {
	var employees []models.Employees
	result := s.DB.Find(&employees)
	if result.Error != nil {
		s.Log.Errorw("Failed to list employees", "error", result.Error)
		res := models.GetApiResponse("api.server.list.error", "FAILED", result.Error.Error())
		return ctx.Status(500).JSON(res)
	}
	s.Log.Infow("Listed employees", "count", len(employees))
	res := models.GetApiResponse("api.server.list.employees", "OK", employees)
	return ctx.JSON(res)
}

func (s *Service) AddEmployee(ctx *fiber.Ctx) error {
	request := new(models.ApiRequest)
	if err := ctx.BodyParser(request); err != nil {
		s.Log.Warnw("Invalid request body for add employee", "error", err)
		res := models.GetApiResponse("api.add", "ERROR", ctx.Status(400).JSON(fiber.Map{"error": err.Error()}))
		return ctx.JSON(res)
	}
	var emp models.Employees
	if err := json.Unmarshal(request.Request, &emp); err != nil {
		s.Log.Warnw("Failed to unmarshal employee data", "error", err)
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if emp.ID == "" {
		emp.ID = uuid.New().String()
	}
	emp.VaultEntityID = emp.ID + "pii"

	// Encrypt Fields
	var err error
	emp.Email, err = s.encryptField(emp.Email, emp.VaultEntityID)
	if err != nil {
		s.Log.Errorw("Failed to encrypt email", "employee_id", emp.ID, "error", err)
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to encrypt email: " + err.Error()})
	}
	emp.Mobile, err = s.encryptField(emp.Mobile, emp.VaultEntityID)
	if err != nil {
		s.Log.Errorw("Failed to encrypt mobile", "employee_id", emp.ID, "error", err)
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to encrypt mobile: " + err.Error()})
	}

	result := s.DB.Create(&emp)
	if result.Error != nil {
		s.Log.Errorw("Failed to create employee", "employee_id", emp.ID, "error", result.Error)
		res := models.GetApiResponse("api.add", "ERROR", result.Error.Error())
		return ctx.Status(500).JSON(res)
	}
	s.Log.Infow("Employee created", "employee_id", emp.ID)
	res := models.GetApiResponse("api.add.employee", "OK", "Inserted Record Successfully")
	return ctx.JSON(res)
}

func (s *Service) GetEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	s.Log.Debugw("Get employee request", "employee_id", id)

	var emp models.Employees
	result := s.DB.First(&emp, "id = ?", id)
	if result.Error != nil {
		s.Log.Warnw("Employee not found", "employee_id", id)
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Decrypt Fields
	var err error
	emp.Email, err = s.decryptField(emp.Email, emp.VaultEntityID)
	if err != nil {
		s.Log.Errorw("Failed to decrypt email", "employee_id", id, "error", err)
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to decrypt email: " + err.Error()})
	}
	emp.Mobile, err = s.decryptField(emp.Mobile, emp.VaultEntityID)
	if err != nil {
		s.Log.Errorw("Failed to decrypt mobile", "employee_id", id, "error", err)
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to decrypt mobile: " + err.Error()})
	}

	s.Log.Infow("Employee retrieved", "employee_id", id)
	res := models.GetApiResponse("api.get.employee", "OK", emp)
	return ctx.JSON(res)
}

func (s *Service) DeleteEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	s.Log.Debugw("Delete employee request", "employee_id", id)

	var emp models.Employees
	result := s.DB.Delete(&emp, "id = ?", id)
	if result.RowsAffected == 0 {
		s.Log.Warnw("Employee not found for deletion", "employee_id", id)
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	if result.Error != nil {
		s.Log.Errorw("Failed to delete employee", "employee_id", id, "error", result.Error)
		return ctx.Status(500).JSON(fiber.Map{"error": "Unable to delete user"})
	}
	s.Log.Infow("Employee deleted", "employee_id", id)
	res := models.GetApiResponse("api.get.employee", "OK", "Record Deleted Successfully")
	return ctx.JSON(res)
}

func (s *Service) UpdateEmployee(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	s.Log.Debugw("Update employee request", "employee_id", id)

	var emp models.Employees
	result := s.DB.First(&emp, "id = ?", id)
	if result.Error != nil {
		s.Log.Warnw("Employee not found for update", "employee_id", id)
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	request := new(models.ApiRequest)
	if err := ctx.BodyParser(request); err != nil {
		s.Log.Warnw("Invalid request body for update employee", "employee_id", id, "error", err)
		res := models.GetApiResponse("api.add", "ERROR", ctx.Status(400).JSON(fiber.Map{"error": err.Error()}))
		return ctx.JSON(res)
	}
	var employee models.Employees
	if err := json.Unmarshal(request.Request, &employee); err != nil {
		s.Log.Warnw("Failed to unmarshal employee update data", "employee_id", id, "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	emp.FirstName = employee.FirstName
	emp.LastName = employee.LastName

	// Encrypt Fields if they are being updated
	encryptedEmail, err := s.encryptField(employee.Email, emp.VaultEntityID)
	if err != nil {
		s.Log.Errorw("Failed to encrypt email on update", "employee_id", id, "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to encrypt email"})
	}
	emp.Email = encryptedEmail

	encryptedMobile, err := s.encryptField(employee.Mobile, emp.VaultEntityID)
	if err != nil {
		s.Log.Errorw("Failed to encrypt mobile on update", "employee_id", id, "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to encrypt mobile"})
	}
	emp.Mobile = encryptedMobile

	emp.Salary = employee.Salary
	emp.VaultEntityID = emp.ID + "pii"

	s.DB.Save(&emp)
	s.Log.Infow("Employee updated", "employee_id", id)
	res := models.GetApiResponse("api.get.employee", "OK", "Record Updated Successfully")
	return ctx.JSON(res)
}

func (s *Service) callKeycloakTokenEndpoint(form url.Values) (map[string]interface{}, int, error) {
	url := fmt.Sprintf(
		"http://%s:%s/realms/%s/protocol/openid-connect/token",
		s.Config.KCHost,
		s.Config.KCPort,
		s.Config.KCRealm,
	)

	resp, err := http.PostForm(url, form)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("failed to reach keycloak: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("failed reading keycloak response: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("failed parsing keycloak response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return data, resp.StatusCode, fmt.Errorf("keycloak returned error: %s", string(body))
	}

	return data, resp.StatusCode, nil
}

func (s *Service) LoginHandler(c *fiber.Ctx) error {

	var req models.LoginRequest[models.Login]

	if err := c.BodyParser(&req); err != nil {
		s.Log.Warnw("Invalid login request body", "error", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	form := url.Values{}
	form.Set("client_id", "employee-api")
	form.Set("grant_type", "password")
	form.Set("username", req.Request.Username)
	form.Set("password", req.Request.Password)

	data, status, err := s.callKeycloakTokenEndpoint(form)

	if err != nil {
		s.Log.Errorw("Keycloak login failed",
			"error", err,
			"status", status,
		)

		return c.Status(status).JSON(data)
	}

	s.Log.Infow("Login successful", "status", status)

	return c.Status(status).JSON(data)
}

func (s *Service) RefreshHandler(c *fiber.Ctx) error {

	var req models.LoginRequest[models.Refresh]

	if err := c.BodyParser(&req); err != nil {
		s.Log.Warnw("Invalid refresh request body", "error", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Request.RefreshToken == "" {
		s.Log.Warn("Refresh token missing")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh_token is required",
		})
	}

	form := url.Values{}
	form.Set("client_id", "employee-api")
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", req.Request.RefreshToken)

	data, status, err := s.callKeycloakTokenEndpoint(form)

	if err != nil {
		s.Log.Errorw("Keycloak refresh failed",
			"error", err,
			"status", status,
		)

		return c.Status(status).JSON(data)
	}

	s.Log.Infow("Token refresh successful", "status", status)

	return c.Status(status).JSON(data)
}
