package routes

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"github.com/colcrunch/avecalc_backend/database"
	"github.com/colcrunch/avecalc_backend/models"
	"github.com/gofiber/fiber/v2"
)

func generateRef() string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, 3)
	rand.Read(b)
	s := hex.EncodeToString(b)
	return s
}

func createContract(c *fiber.Ctx) error {
	// Generate a ref for the contract
	ref := generateRef()

	var co models.Contract

	if err := c.BodyParser(&co); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	co.Ref = ref
	if co.Status > 0 {
		// Contract should never be created with a non-0 status
		co.Status = 0
	}

	database.Db.Create(&co)

	return c.Status(200).JSON(co)
}

func getContracts(c *fiber.Ctx) error {
	contracts := []models.Contract{}

	database.Db.Find(&contracts)

	return c.Status(200).JSON(contracts)
}

func findContractById(id int, co *models.Contract) error {
	database.Db.Find(&co, "id = ?", id)
	if co.ID == 0 {
		return errors.New("contract not found")
	}
	return nil
}

func updateContract(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	type ContractStatus struct {
		Status int8 `json:"status"`
	}

	var co models.Contract

	if err != nil {
		return c.Status(400).JSON(":id must be int")
	}

	if err := findContractById(id, &co); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	var uData ContractStatus
	if err := c.BodyParser(&uData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if uData.Status > 2 || uData.Status < 0 {
		return c.Status(400).JSON("invalid status value")
	}

	co.Status = uData.Status

	database.Db.Save(&co)

	return c.Status(200).JSON(co)
}

func getContract(c *fiber.Ctx) error {
	var co models.Contract

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(":id must be int")
	}
	if err := findContractById(id, &co); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).JSON(co)
}

func deleteContract(c *fiber.Ctx) error {
	var co models.Contract

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(":id must be int")
	}
	if err := findContractById(id, &co); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	if err := database.Db.Delete(&co).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("successfully deleted contract")
}

func ContractRoutes(app *fiber.App) {
	app.Post("/contract", createContract)
	app.Get("/contract/:id", getContract)
	app.Delete("/contract/:id", deleteContract)
	app.Put("/contract/:id", updateContract)
	app.Get("/contracts", getContracts)
}
