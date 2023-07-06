package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("JWT auth")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unathorized")
	}

	if err := parseToken(token); err != nil {
		return err
	}

	fmt.Println("token: ", token)

	return nil
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unathorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("secret: ", secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse jwt token", err)
		return fmt.Errorf("unathorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}

	return fmt.Errorf("unathorized")
}
