package handlers

import (
	"complexity-analyzer/services"

	"github.com/gofiber/fiber/v2"
)

type AnalyzeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type AnalyzeResponse struct {
	Complexity      string             `json:"complexity"`
	Confidence      float64            `json:"confidence"`
	StaticAnalysis  []string           `json:"staticAnalysis"`
	PerformanceData []PerformancePoint `json:"performanceData"`
	Error           string             `json:"error,omitempty"`
}

type PerformancePoint struct {
	Size int     `json:"size"`
	Time float64 `json:"time"`
}

func AnalyzeCode(c *fiber.Ctx) error {
	var req AnalyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Code == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Code is required",
		})
	}

	result, err := services.AnalyzeWithGemini(req.Code, req.Language)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
