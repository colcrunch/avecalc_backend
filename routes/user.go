package routes

import (
	"errors"
	"log"

	"github.com/colcrunch/avecalc_backend/database"
	"github.com/colcrunch/avecalc_backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserSerializer struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func serializeUser(modelUser models.User) UserSerializer {
	return UserSerializer{
		ID:       modelUser.ID,
		UserName: modelUser.UserName,
		Email:    modelUser.Email,
		Admin:    modelUser.Admin,
	}
}

func createUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	user.PasswordHash = getHash([]byte(user.PasswordHash))

	database.Db.Create(&user)

	responseUser := serializeUser(user)

	return c.Status(201).JSON(responseUser)
}

func findUserById(id int, user *models.User) error {
	database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func getUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON(":id must be int.")
	}

	if err := findUserById(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	responseUser := serializeUser(user)

	return c.Status(200).JSON(responseUser)
}

func getUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Db.Find(&users)

	responseUsers := []UserSerializer{}

	for _, user := range users {
		responseUser := serializeUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func updateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON(":id should be int")
	}

	var user models.User

	if err := findUserById(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	type UpdateUser struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Admin    string `json:"admin"`
	}

	var uData UpdateUser

	if err := c.BodyParser(&uData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if uData.UserName != "" {
		user.UserName = uData.UserName
	}
	if uData.Email != "" {
		user.Email = uData.Email
	}
	if uData.Password != "" {
		user.PasswordHash = getHash([]byte(user.PasswordHash))
	}
	if uData.Admin != "" {
		if uData.Admin == "true" {
			user.Admin = true
		} else {
			user.Admin = false
		}
	}

	database.Db.Save(&user)
	responseUser := serializeUser(user)

	return c.Status(200).JSON(responseUser)
}

func deleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON(":id must be int.")
	}

	if err := findUserById(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	if err := database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully deleted user.")
}

func UserRoutes(app *fiber.App) {
	app.Get("/users", getUsers)
	app.Post("/user", createUser)
	app.Get("/user/:id", getUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)
}
