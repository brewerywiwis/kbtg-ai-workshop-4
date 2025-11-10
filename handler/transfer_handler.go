package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"workshop4-backend/service"
)

type TransferCreateRequest struct {
	FromUserID int     `json:"fromUserId" validate:"required,min=1"`
	ToUserID   int     `json:"toUserId" validate:"required,min=1"`
	Amount     int     `json:"amount" validate:"required,min=1"`
	Note       *string `json:"note,omitempty"`
}

type TransferCreateResponse struct {
	Transfer interface{} `json:"transfer"`
}

type TransferGetResponse struct {
	Transfer interface{} `json:"transfer"`
}

type TransferListResponse struct {
	Data     interface{} `json:"data"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type TransferHandler struct {
	service *service.TransferService
}

func NewTransferHandler(service *service.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}

func (h *TransferHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/transfers", h.CreateTransfer)
	app.Get("/transfers", h.GetTransfers)
	app.Get("/transfers/:id", h.GetTransferByID)
}

func (h *TransferHandler) CreateTransfer(c *fiber.Ctx) error {
	var req TransferCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
	}

	// Validate required fields
	if req.FromUserID <= 0 {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "fromUserId must be greater than 0",
		})
	}
	if req.ToUserID <= 0 {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "toUserId must be greater than 0",
		})
	}
	if req.Amount <= 0 {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "amount must be greater than 0",
		})
	}

	transfer, err := h.service.CreateTransfer(req.FromUserID, req.ToUserID, req.Amount, req.Note)
	if err != nil {
		switch err {
		case service.ErrSelfTransfer:
			return c.Status(422).JSON(ErrorResponse{
				Error:   "SELF_TRANSFER",
				Message: "Cannot transfer to yourself",
			})
		case service.ErrInsufficientBalance:
			return c.Status(409).JSON(ErrorResponse{
				Error:   "INSUFFICIENT_BALANCE",
				Message: "Insufficient balance",
			})
		case service.ErrUserNotFound:
			return c.Status(400).JSON(ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
		default:
			return c.Status(500).JSON(ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to create transfer",
			})
		}
	}

	// Set idempotency key header
	c.Set("Idempotency-Key", transfer.IdempotencyKey)

	return c.Status(201).JSON(TransferCreateResponse{
		Transfer: transfer,
	})
}

func (h *TransferHandler) GetTransferByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Transfer ID is required",
		})
	}

	transfer, err := h.service.GetTransferByIdempotencyKey(id)
	if err != nil {
		if err == service.ErrTransferNotFound {
			return c.Status(404).JSON(ErrorResponse{
				Error:   "NOT_FOUND",
				Message: "Transfer not found",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to get transfer",
		})
	}

	return c.JSON(TransferGetResponse{
		Transfer: transfer,
	})
}

func (h *TransferHandler) GetTransfers(c *fiber.Ctx) error {
	// Get userID from query parameter
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "userId query parameter is required",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return c.Status(400).JSON(ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "userId must be a valid positive integer",
		})
	}

	// Get pagination parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 200 {
			pageSize = ps
		}
	}

	transfers, total, err := h.service.GetTransfersByUserID(userID, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to get transfers",
		})
	}

	return c.JSON(TransferListResponse{
		Data:     transfers,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}