package helper

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"instagram-bot-live/internal/domain"
	"strings"
	"time"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(p string) (string, error) {

	if len(p) < 1 {
		return "", errors.New("Password is too short!")
	}

	hashP, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		// log actual error and report to loggin tool
		return "", errors.New("Password hash is failed")
	}

	return string(hashP), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {

	if id == 0 || email == "" || role == "" {
		return "", errors.New("Invalid parameter to generate the token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(a.Secret))

	if err != nil {
		// log the erros
		return "", errors.New("Token signing failed")
	}

	return tokenStr, nil
}

func (a Auth) VerifyPassword(pP string, hP string) error {
	if len(pP) < 1 {
		return errors.New("Password is too short!")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hP), []byte(pP))
	if err != nil {
		return errors.New("Password or email is invalid")
	}

	return nil
}

func (a Auth) verifyToken(t string) (domain.User, error) {
	tokenArr := strings.Split(t, " ")
	if len(tokenArr) != 2 {
		return domain.User{}, errors.New("Token is invalid")
	}

	tokenStr := tokenArr[1]
	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("Token is invalid")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(a.Secret), nil
	})
	if err != nil {
		return domain.User{}, errors.New("Token is invalid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("Token is expired")
		}

		user := domain.User{}
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.UserType = claims["role"].(string)
		return user, nil
	}

	return domain.User{}, errors.New("Token verification failed")
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {

	authHeader := ctx.GetReqHeaders()["Authorization"]
	if authHeader == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Authorization failed",
			"reason":  "Token not informed",
		})
	}

	user, err := a.verifyToken(authHeader[0])
	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
		"message": "Authorization failed",
		"reason":  err,
	})
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {

	user := ctx.Locals("user").(domain.User)

	return user
}

func (a Auth) GenerateCode() (int, error) {
	return RandomNumbers(6)
}

func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {

	authHeader := ctx.GetReqHeaders()["Authorization"]
	if authHeader == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Authorization failed",
			"reason":  "Token not informed",
		})
	}

	user, err := a.verifyToken(authHeader[0])
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Authorization failed",
			"reason":  err,
		})
	}

	if user.ID > 0 && user.UserType == domain.SELLER {
		ctx.Locals("user", user)
		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
		"message": "Authorization failed",
		"reason":  errors.New("please join seller program to manage products"),
	})
}
