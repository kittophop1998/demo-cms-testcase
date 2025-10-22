package handlers

import (
	"demo-notion-api/config"
	"demo-notion-api/models"
	"demo-notion-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotionHandler struct {
	notionService *services.NotionService
}

func NewNotionHandler(cfg *config.Config) *NotionHandler {
	return &NotionHandler{
		notionService: services.NewNotionService(cfg),
	}
}

// SearchTestCases godoc
// @Summary Search for test cases
// @Description Search for test cases that start with TC_ prefix from Notion
// @Tags testcases
// @Accept json
// @Produce json
// @Success 200 {array} models.TestCaseResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/test-cases [get]
func (h *NotionHandler) SearchTestCases(c *gin.Context) {
	testCases, err := h.notionService.SearchTestCases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to search test cases",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    testCases,
		Message: "Test cases retrieved successfully",
	})
}

// GetTestCaseBlocks godoc
// @Summary Get blocks for a specific test case
// @Description Get all blocks from a test case page by test case key
// @Tags testcases
// @Accept json
// @Produce json
// @Param testCaseKey path string true "Test Case Key (e.g., 01001)"
// @Param type query string false "Block type filter (e.g., table)"
// @Success 200 {array} models.BlockResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/test-cases/{testCaseKey}/blocks [get]
func (h *NotionHandler) GetTestCaseBlocks(c *gin.Context) {
	testCaseKey := c.Param("testCaseKey")
	if testCaseKey == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing test case key",
			Message: "Test case key is required",
		})
		return
	}

	// Find the test case by key
	testCase, err := h.notionService.GetTestCaseByKey(testCaseKey)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Test case not found",
			Message: err.Error(),
		})
		return
	}

	// Check if user wants only table blocks
	blockType := c.Query("type")

	var blocks []models.BlockResponse
	if blockType == "table" {
		blocks, err = h.notionService.GetTableBlocks(testCase.PageID)
	} else {
		blocks, err = h.notionService.GetPageBlocks(testCase.PageID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get blocks",
			Message: err.Error(),
		})
		return
	}

	response := TestCaseBlocksResponse{
		TestCase: *testCase,
		Blocks:   blocks,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
		Message: "Blocks retrieved successfully",
	})
}

// GetBlockDetails godoc
// @Summary Get detailed information about a specific block
// @Description Get detailed information about a block by its ID
// @Tags blocks
// @Accept json
// @Produce json
// @Param blockId path string true "Block ID"
// @Success 200 {object} models.BlockResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/blocks/{blockId} [get]
func (h *NotionHandler) GetBlockDetails(c *gin.Context) {
	blockID := c.Param("blockId")
	if blockID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing block ID",
			Message: "Block ID is required",
		})
		return
	}

	block, err := h.notionService.GetBlockDetails(blockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get block details",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    block,
		Message: "Block details retrieved successfully",
	})
}

// Response structures for API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type TestCaseBlocksResponse struct {
	TestCase models.TestCaseResponse `json:"test_case"`
	Blocks   []models.BlockResponse  `json:"blocks"`
}
