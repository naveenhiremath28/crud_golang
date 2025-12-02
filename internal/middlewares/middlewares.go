package middlewares

import (
	"log"
	"practise/go_fiber/internal/models"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func KeycloakAuth(jwksURL string) fiber.Handler {
	// Load JWKS once
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		// DO NOT IGNORE THIS ERROR
		log.Fatalf("Failed to load JWKS from Keycloak: %v", err)
	}

	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			res := models.GetApiResponse("api", "ERROR", "Missing token")
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		// Extract token
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		// Validate token
		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			res := models.GetApiResponse("api", "ERROR", err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		return c.Next()
	}
}
