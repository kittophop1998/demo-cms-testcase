package main

import (
	"demo-notion-api/config"
	"demo-notion-api/services"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// This is a simple test script to demonstrate the Notion service functionality
// Run with: go run examples/test_notion.go
// Make sure to set your NOTION_API_KEY environment variable before running

func main() {
	// Load configuration
	cfg := config.Load()

	if cfg.NotionAPIKey == "" {
		log.Fatal("NOTION_API_KEY environment variable is required")
	}

	// Create notion service
	notionService := services.NewNotionService(cfg)

	fmt.Println("🔍 Searching for test cases...")

	// Search for test cases
	testCases, err := notionService.SearchTestCases()
	if err != nil {
		log.Fatalf("Error searching test cases: %v", err)
	}

	fmt.Printf("Found %d test cases:\n", len(testCases))

	for i, tc := range testCases {
		fmt.Printf("%d. TC_%s: %s (Status: %s)\n", i+1, tc.TestCaseKey, tc.Title, tc.Status)

		// If this is the first test case, get its blocks
		if i == 0 {
			fmt.Printf("\n📄 Getting blocks for test case TC_%s...\n", tc.TestCaseKey)

			blocks, err := notionService.GetPageBlocks(tc.PageID)
			if err != nil {
				fmt.Printf("Error getting blocks: %v\n", err)
				continue
			}

			fmt.Printf("Found %d blocks:\n", len(blocks))
			for j, block := range blocks {
				fmt.Printf("  %d. Block Type: %s, ID: %s", j+1, block.Type, block.BlockID)
				if block.Type == "table" && block.TableInfo != nil {
					fmt.Printf(" (Table: %dx%d)", block.TableInfo.TableWidth, block.TableInfo.TableWidth)
				}
				fmt.Println()
			}

			// Get table blocks specifically
			fmt.Printf("\n📊 Getting table blocks for test case TC_%s...\n", tc.TestCaseKey)
			tableBlocks, err := notionService.GetTableBlocks(tc.PageID)
			if err != nil {
				fmt.Printf("Error getting table blocks: %v\n", err)
				continue
			}

			fmt.Printf("Found %d table blocks:\n", len(tableBlocks))
			for j, block := range tableBlocks {
				fmt.Printf("  %d. Table Block ID: %s", j+1, block.BlockID)
				if block.TableInfo != nil {
					fmt.Printf(" (Size: %dx?, Headers: Col=%t, Row=%t)",
						block.TableInfo.TableWidth,
						block.TableInfo.HasColumnHeader,
						block.TableInfo.HasRowHeader)
				}
				fmt.Println()
			}
		}
	}

	if len(testCases) > 0 {
		// Pretty print the first test case as JSON
		fmt.Println("\n📋 Sample test case data (JSON):")
		jsonData, err := json.MarshalIndent(testCases[0], "", "  ")
		if err == nil {
			fmt.Println(string(jsonData))
		}
	}
}

func init() {
	// Load .env file if it exists (optional, for testing)
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println("💡 Tip: .env file found. Make sure to set NOTION_API_KEY in it or as environment variable")
	}
}
