package handler

import (
	"strconv"
	"time"

	"workshop4-backend/domain"
	"workshop4-backend/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/users", h.GetAllUsers)
	app.Get("/users/:id", h.GetUserByID)
	app.Post("/users", h.CreateUser)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.JSON(users)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var newUser domain.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	if err := h.service.CreateUser(&newUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}
	return c.Status(201).JSON(newUser)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	var updateUser domain.User
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	updateUser.ID = id
	updateUser.UpdatedAt = time.Now()
	if err := h.service.UpdateUser(&updateUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}
	return c.JSON(updateUser)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
