package api

import (
	"fmt"
	"os"
	"time"

	"github.com/boyanivskyy/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrorUnathorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			fmt.Println("invalid token(JWTAuthentication)")
			return ErrorUnathorized()
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			fmt.Println("token expired")
			return ErrorUnathorized()
		}

		userId := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userId)
		if err != nil {
			fmt.Println("no user related to token")
			return ErrorUnathorized()
		}

		// set the current authenticated user to the context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrorUnathorized()
		}
		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse jwt token", err)
		return nil, ErrorUnathorized()
	}

	if !token.Valid {
		fmt.Println("invalid token(validateToken)", token)
		return nil, ErrorUnathorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrorUnathorized()
	}

	return claims, nil
}
