package models

import "time"

// NotionSearchRequest represents the search request payload
type NotionSearchRequest struct {
	Query  string             `json:"query"`
	Filter NotionSearchFilter `json:"filter"`
	Sort   NotionSearchSort   `json:"sort"`
}

type NotionSearchFilter struct {
	Value    string `json:"value"`
	Property string `json:"property"`
}

type NotionSearchSort struct {
	Direction string `json:"direction"`
	Timestamp string `json:"timestamp"`
}

// NotionSearchResponse represents the search response
type NotionSearchResponse struct {
	Object     string       `json:"object"`
	Results    []NotionPage `json:"results"`
	NextCursor string       `json:"next_cursor"`
	HasMore    bool         `json:"has_more"`
}

// NotionPage represents a page in Notion
type NotionPage struct {
	Object         string                 `json:"object"`
	ID             string                 `json:"id"`
	CreatedTime    time.Time              `json:"created_time"`
	LastEditedTime time.Time              `json:"last_edited_time"`
	CreatedBy      NotionUser             `json:"created_by"`
	LastEditedBy   NotionUser             `json:"last_edited_by"`
	Cover          interface{}            `json:"cover"`
	Icon           interface{}            `json:"icon"`
	Parent         NotionParent           `json:"parent"`
	Archived       bool                   `json:"archived"`
	InTrash        bool                   `json:"in_trash"`
	IsLocked       bool                   `json:"is_locked"`
	Properties     map[string]interface{} `json:"properties"`
	URL            string                 `json:"url"`
	PublicURL      interface{}            `json:"public_url"`
}

type NotionUser struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type NotionParent struct {
	Type         string `json:"type"`
	DataSourceID string `json:"data_source_id"`
	DatabaseID   string `json:"database_id"`
}

// NotionProperty represents different types of properties
type NotionTitleProperty struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Title []TitleText `json:"title"`
}

type TitleText struct {
	Type        string      `json:"type"`
	Text        TextContent `json:"text"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        interface{} `json:"href"`
}

type TextContent struct {
	Content string      `json:"content"`
	Link    interface{} `json:"link"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

// NotionBlocksResponse represents the blocks response
type NotionBlocksResponse struct {
	Object     string        `json:"object"`
	Results    []NotionBlock `json:"results"`
	NextCursor string        `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}

// NotionBlock represents a block in Notion
type NotionBlock struct {
	Object         string            `json:"object"`
	ID             string            `json:"id"`
	Parent         NotionBlockParent `json:"parent"`
	CreatedTime    time.Time         `json:"created_time"`
	LastEditedTime time.Time         `json:"last_edited_time"`
	CreatedBy      NotionUser        `json:"created_by"`
	LastEditedBy   NotionUser        `json:"last_edited_by"`
	HasChildren    bool              `json:"has_children"`
	Archived       bool              `json:"archived"`
	InTrash        bool              `json:"in_trash"`
	Type           string            `json:"type"`
	Table          *NotionTable      `json:"table,omitempty"`
	TableRow       *NotionTableRow   `json:"table_row,omitempty"`
	Paragraph      *NotionParagraph  `json:"paragraph,omitempty"`
	Heading2       *NotionHeading    `json:"heading_2,omitempty"`
}

type NotionBlockParent struct {
	Type   string `json:"type"`
	PageID string `json:"page_id"`
}

type NotionTable struct {
	TableWidth      int  `json:"table_width"`
	HasColumnHeader bool `json:"has_column_header"`
	HasRowHeader    bool `json:"has_row_header"`
}

type NotionParagraph struct {
	RichText []RichText `json:"rich_text"`
}

type NotionHeading struct {
	RichText []RichText `json:"rich_text"`
}

type RichText struct {
	Type        string      `json:"type"`
	Text        TextContent `json:"text"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        interface{} `json:"href"`
}

// TestCaseResponse represents our custom response for test cases
type TestCaseResponse struct {
	TestCaseKey string    `json:"test_case_key"`
	PageID      string    `json:"page_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	TestDate    string    `json:"test_date"`
	URL         string    `json:"url"`
	LastEdited  time.Time `json:"last_edited"`
}

// BlockResponse represents our custom response for blocks
type BlockResponse struct {
	BlockID     string     `json:"block_id"`
	Type        string     `json:"type"`
	HasChildren bool       `json:"has_children"`
	Content     string     `json:"content,omitempty"`
	TableInfo   *TableInfo `json:"table_info,omitempty"`
}

type TableInfo struct {
	TableWidth      int        `json:"table_width"`
	HasColumnHeader bool       `json:"has_column_header"`
	HasRowHeader    bool       `json:"has_row_header"`
	Rows            []TableRow `json:"rows,omitempty"`
}

// TableRow represents a row in a table
type TableRow struct {
	Cells []string `json:"cells"`
}

// Detailed test case response with table data
type DetailedTestCaseResponse struct {
	TestCaseKey string          `json:"test_case_key"`
	PageID      string          `json:"page_id"`
	Title       string          `json:"title"`
	Status      string          `json:"status"`
	TestDate    string          `json:"test_date"`
	URL         string          `json:"url"`
	LastEdited  time.Time       `json:"last_edited"`
	Tables      []TableWithData `json:"tables,omitempty"`
}

type TableWithData struct {
	BlockID         string     `json:"block_id"`
	TableWidth      int        `json:"table_width"`
	HasColumnHeader bool       `json:"has_column_header"`
	HasRowHeader    bool       `json:"has_row_header"`
	Rows            []TableRow `json:"rows"`
}

// Table row block structure from Notion API
type NotionTableRow struct {
	Cells [][]RichText `json:"cells"`
}
