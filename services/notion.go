package services

import (
	"bytes"
	"demo-notion-api/config"
	"demo-notion-api/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type NotionService struct {
	config *config.Config
	client *http.Client
}

func NewNotionService(cfg *config.Config) *NotionService {
	return &NotionService{
		config: cfg,
		client: &http.Client{},
	}
}

// SearchTestCases searches for pages with "External tasks" query and extracts test cases
func (s *NotionService) SearchTestCases() ([]models.TestCaseResponse, error) {
	searchReq := models.NotionSearchRequest{
		Filter: models.NotionSearchFilter{
			Value:    "page",
			Property: "object",
		},
		Sort: models.NotionSearchSort{
			Direction: "ascending",
			Timestamp: "last_edited_time",
		},
	}

	reqBody, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.NotionAPIURL+"/search", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var searchResp models.NotionSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("searchResp.Results", searchResp.Results)
	return s.extractTestCases(searchResp.Results), nil
}

// GetTestCaseByKey finds a test case by its key (e.g., "01001" from "TC_01001")
func (s *NotionService) GetTestCaseByKey(testCaseKey string) (*models.TestCaseResponse, error) {
	testCases, err := s.SearchTestCases()
	if err != nil {
		return nil, err
	}

	for _, tc := range testCases {
		if tc.TestCaseKey == testCaseKey {
			return &tc, nil
		}
	}

	return nil, fmt.Errorf("test case with key %s not found", testCaseKey)
}

// GetPageBlocks retrieves all blocks from a page
func (s *NotionService) GetPageBlocks(pageID string) ([]models.BlockResponse, error) {
	url := fmt.Sprintf("%s/blocks/%s/children", s.config.NotionAPIURL, pageID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var blocksResp models.NotionBlocksResponse
	if err := json.NewDecoder(resp.Body).Decode(&blocksResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.convertToBlockResponses(blocksResp.Results), nil
}

// GetBlockDetails retrieves detailed information about a specific block
func (s *NotionService) GetBlockDetails(blockID string) (*models.BlockResponse, error) {
	url := fmt.Sprintf("%s/blocks/%s", s.config.NotionAPIURL, blockID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var block models.NotionBlock
	if err := json.NewDecoder(resp.Body).Decode(&block); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.convertToBlockResponse(block), nil
}

// GetTableBlocks filters blocks to return only table type blocks
func (s *NotionService) GetTableBlocks(pageID string) ([]models.BlockResponse, error) {
	blocks, err := s.GetPageBlocks(pageID)
	if err != nil {
		return nil, err
	}

	var tableBlocks []models.BlockResponse
	for _, block := range blocks {
		if block.Type == "table" {
			tableBlocks = append(tableBlocks, block)
		}
	}

	return tableBlocks, nil
}

// GetDetailedTestCases searches for test cases and includes their table data
func (s *NotionService) GetDetailedTestCases() ([]models.DetailedTestCaseResponse, error) {
	// First get all test cases
	testCases, err := s.SearchTestCases()
	if err != nil {
		return nil, fmt.Errorf("failed to search test cases: %w", err)
	}

	var detailedTestCases []models.DetailedTestCaseResponse

	for _, tc := range testCases {
		detailed := models.DetailedTestCaseResponse{
			TestCaseKey: tc.TestCaseKey,
			PageID:      tc.PageID,
			Title:       tc.Title,
			Status:      tc.Status,
			TestDate:    tc.TestDate,
			URL:         tc.URL,
			LastEdited:  tc.LastEdited,
		}

		// Get table blocks for this test case
		tableBlocks, err := s.GetTableBlocks(tc.PageID)
		if err != nil {
			// Log error but continue with other test cases
			fmt.Printf("Warning: failed to get table blocks for test case %s: %v\n", tc.TestCaseKey, err)
			detailed.Tables = []models.TableWithData{}
		} else {
			// Get table data for each table block
			for _, tableBlock := range tableBlocks {
				tableData, err := s.GetTableData(tableBlock.BlockID)
				if err != nil {
					fmt.Printf("Warning: failed to get table data for block %s: %v\n", tableBlock.BlockID, err)
					continue
				}
				detailed.Tables = append(detailed.Tables, *tableData)
			}
		}

		detailedTestCases = append(detailedTestCases, detailed)
	}

	return detailedTestCases, nil
}

// GetTableData retrieves table data including all rows
func (s *NotionService) GetTableData(tableBlockID string) (*models.TableWithData, error) {
	// First get the table block info
	tableBlock, err := s.GetBlockDetails(tableBlockID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table block details: %w", err)
	}

	if tableBlock.Type != "table" || tableBlock.TableInfo == nil {
		return nil, fmt.Errorf("block %s is not a table", tableBlockID)
	}

	// Get table rows (children of the table block)
	url := fmt.Sprintf("%s/blocks/%s/children", s.config.NotionAPIURL, tableBlockID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var blocksResp models.NotionBlocksResponse
	if err := json.NewDecoder(resp.Body).Decode(&blocksResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert table rows
	var rows []models.TableRow
	for _, block := range blocksResp.Results {
		if block.Type == "table_row" && block.TableRow != nil {
			var cells []string
			for _, cellArray := range block.TableRow.Cells {
				cellContent := s.extractRichTextContent(cellArray)
				cells = append(cells, cellContent)
			}
			rows = append(rows, models.TableRow{Cells: cells})
		}
	}

	tableData := &models.TableWithData{
		BlockID:         tableBlockID,
		TableWidth:      tableBlock.TableInfo.TableWidth,
		HasColumnHeader: tableBlock.TableInfo.HasColumnHeader,
		HasRowHeader:    tableBlock.TableInfo.HasRowHeader,
		Rows:            rows,
	}

	return tableData, nil
}

// Helper methods

func (s *NotionService) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+s.config.NotionAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", s.config.NotionAPIVersion)
}

func (s *NotionService) extractTestCases(pages []models.NotionPage) []models.TestCaseResponse {
	var testCases []models.TestCaseResponse

	// Regex to match TC_ followed by numbers
	tcRegex := regexp.MustCompile(`TC_(\d+)`)
	fmt.Println("pages", pages)

	for _, page := range pages {
		// Extract Test Case Name property
		if testCaseNameProp, exists := page.Properties["Test Case Name"]; exists {
			if titleProp, ok := testCaseNameProp.(map[string]interface{}); ok {
				if titleArray, exists := titleProp["title"].([]interface{}); exists && len(titleArray) > 0 {
					if firstTitle, ok := titleArray[0].(map[string]interface{}); ok {
						if plainText, exists := firstTitle["plain_text"].(string); exists {
							// Check if it starts with TC_ and extract the key
							if matches := tcRegex.FindStringSubmatch(plainText); matches != nil {
								testCaseKey := matches[1] // Extract the number part

								// Extract other properties
								status := s.extractStatus(page.Properties)
								testDate := s.extractTestDate(page.Properties)

								testCase := models.TestCaseResponse{
									TestCaseKey: testCaseKey,
									PageID:      page.ID,
									Title:       plainText,
									Status:      status,
									TestDate:    testDate,
									URL:         page.URL,
									LastEdited:  page.LastEditedTime,
								}

								testCases = append(testCases, testCase)
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("Extracted %d test cases\n", len(testCases))

	return testCases
}

func (s *NotionService) extractStatus(properties map[string]interface{}) string {
	if statusProp, exists := properties["Status"]; exists {
		if statusMap, ok := statusProp.(map[string]interface{}); ok {
			if status, exists := statusMap["status"].(map[string]interface{}); exists {
				if name, exists := status["name"].(string); exists {
					return name
				}
			}
		}
	}
	return ""
}

func (s *NotionService) extractTestDate(properties map[string]interface{}) string {
	if dateProp, exists := properties["Test Date"]; exists {
		if dateMap, ok := dateProp.(map[string]interface{}); ok {
			if date, exists := dateMap["date"].(map[string]interface{}); exists {
				if start, exists := date["start"].(string); exists {
					return start
				}
			}
		}
	}
	return ""
}

func (s *NotionService) convertToBlockResponses(blocks []models.NotionBlock) []models.BlockResponse {
	var blockResponses []models.BlockResponse

	for _, block := range blocks {
		blockResponses = append(blockResponses, *s.convertToBlockResponse(block))
	}

	return blockResponses
}

func (s *NotionService) convertToBlockResponse(block models.NotionBlock) *models.BlockResponse {
	blockResp := &models.BlockResponse{
		BlockID:     block.ID,
		Type:        block.Type,
		HasChildren: block.HasChildren,
	}

	// Extract content based on block type
	switch block.Type {
	case "table":
		if block.Table != nil {
			blockResp.TableInfo = &models.TableInfo{
				TableWidth:      block.Table.TableWidth,
				HasColumnHeader: block.Table.HasColumnHeader,
				HasRowHeader:    block.Table.HasRowHeader,
			}
		}
	case "paragraph":
		if block.Paragraph != nil {
			blockResp.Content = s.extractRichTextContent(block.Paragraph.RichText)
		}
	case "heading_2":
		if block.Heading2 != nil {
			blockResp.Content = s.extractRichTextContent(block.Heading2.RichText)
		}
	}

	return blockResp
}

func (s *NotionService) extractRichTextContent(richText []models.RichText) string {
	var content []string
	for _, rt := range richText {
		content = append(content, rt.PlainText)
	}
	return strings.Join(content, "")
}
