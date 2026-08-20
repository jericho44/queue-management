package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Meta struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	RequestID  string          `json:"request_id"`
	Timestamp  string          `json:"timestamp"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	From        int   `json:"from"`
	To          int   `json:"to"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    Meta        `json:"meta"`
}

func getRequestID(c *fiber.Ctx) string {
	reqID := c.Get("X-Request-ID")
	if reqID == "" {
		if val, ok := c.Locals("request_id").(string); ok && val != "" {
			reqID = val
		} else {
			reqID = uuid.New().String()
		}
	}
	c.Set("X-Request-ID", reqID)
	return reqID
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	reqID := getRequestID(c)
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			RequestID: reqID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(statusCode).JSON(resp)
}

func Paginated(c *fiber.Ctx, message string, data interface{}, page, perPage int, total int64) error {
	reqID := getRequestID(c)
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}

	from := (page-1)*perPage + 1
	if total == 0 {
		from = 0
	}
	to := page * perPage
	if int64(to) > total {
		to = int(total)
	}

	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			Pagination: &PaginationMeta{
				CurrentPage: page,
				LastPage:    lastPage,
				PerPage:     perPage,
				Total:       total,
				From:        from,
				To:          to,
				HasNext:     page < lastPage,
				HasPrevious: page > 1,
			},
			RequestID: reqID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func Error(c *fiber.Ctx, statusCode int, message string, errs interface{}) error {
	reqID := getRequestID(c)
	resp := APIResponse{
		Success: false,
		Message: message,
		Errors:  errs,
		Meta: Meta{
			RequestID: reqID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.Status(statusCode).JSON(resp)
}
