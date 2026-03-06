package middlewares

import (
	"fmt"
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/models"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func KeycloakAuth(log *zap.SugaredLogger, cfg *config.Config) fiber.Handler {
	url := fmt.Sprintf("http://%s:%s/realms/%s/protocol/openid-connect/certs", cfg.KCHost, cfg.KCPort, cfg.KCRealm)
	jwks, err := keyfunc.Get(url, keyfunc.Options{})
	if err != nil {
		log.Fatalw("Failed to load JWKS from Keycloak",
			"jwks_url", url,
			"error", err,
		)
	}

	log.Infow("JWKS loaded successfully", "jwks_url", url)

	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			log.Debug("Request missing authorization token")
			res := models.GetApiResponse("api", "ERROR", "Missing token")
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			log.Error("Invalid JWT token", "error", err)
			res := models.GetApiResponse("api", "ERROR", err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_claims", claims)

		return c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user_claims").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.GetApiResponse("api", "ERROR", "Unauthorized"))
		}

		realmAccess, ok := claims["realm_access"].(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.GetApiResponse("api", "ERROR", "Access denied: No roles found"))
		}

		roles, ok := realmAccess["roles"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.GetApiResponse("api", "ERROR", "Access denied: Invalid roles format"))
		}

		roleMap := make(map[string]bool)
		for _, r := range roles {
			if roleStr, ok := r.(string); ok {
				roleMap[roleStr] = true
			}
		}

		for _, allowedRole := range allowedRoles {
			if roleMap[allowedRole] {
				return c.Next()
			}
		}


		return c.Status(fiber.StatusForbidden).JSON(models.GetApiResponse("api", "ERROR", "Access denied: Insufficient privileges"))
	}
}
