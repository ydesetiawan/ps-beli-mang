package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ps-beli-mang/pkg/httphelper"
	"ps-beli-mang/pkg/httphelper/response"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

var jwtKey = []byte("your_secret_key")

func JWTAuthMiddleware(fn echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtToken, err := extractJWTTokenFromHeader(c.Request())
		if err != nil {
			slog.Error("Failed to extract JWT token from header", "error", err)
			writeUnauthorized(c.Response())
			return err
		}

		claims, err := parseJWTToClaims(jwtToken)
		if err != nil {
			slog.Error("Failed to parse JWT token", "error", err)
			writeUnauthorized(c.Response())
			return err
		}

		userId, uidOk := claims["user_id"].(string)
		role, roleOk := claims["role"].(string)
		if !uidOk || !roleOk {
			slog.Error("Invalid claims")
			writeUnauthorized(c.Response())
			return err
		}

		user, err := constructUserInfo(userId, role)
		if err != nil {
			slog.Error("Failed to construct user info", "error", err)
			writeUnauthorized(c.Response())
			return err
		}

		c.Set("user_info", user)

		return fn(c)
	}
}

func writeUnauthorized(rw http.ResponseWriter) {
	httphelper.WriteJSON(
		rw, http.StatusUnauthorized,
		response.WebResponse{
			Status:  http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		},
	)
}

func extractJWTTokenFromHeader(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return "", fmt.Errorf("missing auth token")
	}

	return authToken[len("Bearer "):], nil
}

type Claims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.Claims
}

func GenerateJWT(userId string, role string) (string, error) {
	// Create token
	claims := Claims{
		UserId: userId,
		Role:   role,
		Claims: jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it
	return token.SignedString(jwtKey)
}

func constructUserInfo(userId string, role string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"user_id": userId,
		"role":    role,
	}, nil
}

func parseJWTToClaims(jwtToken string) (jwt.MapClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(jwtToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	// no need to verify 'token' signature since it already validated in authz kong plugin, just parse the token

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid jwt token")
	}

	return claims, nil
}
