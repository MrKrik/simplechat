package jwt

import (
	"auth/iternal/domain/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NewToken creates new JWT token for given user and app.
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Добавляем в токен всю необходимую информацию
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	// Подписываем токен, используя секретный ключ приложения
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetLoginFromToken(tokenString string, secret string) (login string, appID float64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(err)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		login, ok := claims["login"].(string)
		if !ok {
			return "", 0, errors.New("login not found in token claims")
		}
		appID, ok := claims["app_id"].(float64)
		if !ok {
			return "", 0, errors.New("appID not found in token claims")
		}
		return login, appID, nil
	}

	return "", 0, errors.New("invalid token")
}

func NewChatToken(login string, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", errors.New("failed to sign token:")
	}

	return tokenString, nil
}
