package routes

import (
	"errors"
	"time"

	"github.com/colcrunch/avecalc_backend/database"
	"github.com/colcrunch/avecalc_backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func createJwt(user models.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 60).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte("supersecret"))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func signup(c *fiber.Ctx) error {
	var user models.User

	type SignupForm struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	var sData SignupForm

	if err := c.BodyParser(&sData); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if sData.UserName == "" || sData.Password == "" {
		return c.Status(400).JSON("Must provide a username and password.")
	}

	user.UserName = sData.UserName
	user.PasswordHash = getHash([]byte(sData.Password))

	database.Db.Create(&user)
	if user.ID == 0 {
		return c.Status(400).JSON("username not availible")
	}

	token, exp, err := createJwt(user)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{"token": token, "exp": exp, "user": serializeUser(user)})

}

func findUserByName(username string, user *models.User) error {
	database.Db.Find(&user, "username = ?", username)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func login(c *fiber.Ctx) error {
	var user models.User

	type LoginForm struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	var lData LoginForm

	if err := c.BodyParser(&lData); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if lData.Password == "" || lData.UserName == "" {
		return c.Status(400).JSON("username and password required")
	}

	if err := findUserByName(lData.UserName, &user); err != nil {
		return c.Status(404).JSON("could not find a user with those credentials")
	}

	if getHash([]byte(lData.Password)) != user.PasswordHash {
		return c.Status(404).JSON("could not find a user with those credentials")
	}

	// Successful Login
	token, exp, err := createJwt(user)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{"token": token, "exp": exp, "user": serializeUser(user)})
}

func AuthRoutes(app *fiber.App) {
	app.Post("/auth/signup", signup)
	app.Post("/auth/login", login)
}
