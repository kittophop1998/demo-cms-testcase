# Demo Notion API - Go Gin Backend

A Go backend service using Gin framework that integrates with Notion API to manage test cases and blocks.

## Features

- 🔍 Search for test cases with TC_ prefix from Notion database
- 📄 Retrieve page blocks from specific test cases  
- 📊 Filter and retrieve table blocks
- 🎯 Extract test case information (status, dates, etc.)
- 🚀 RESTful API endpoints

## Prerequisites

- Go 1.19 or higher
- Notion API key
- Notion database with test cases

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd demo-notion-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env file with your Notion API key
```

## Configuration

Create a `.env` file based on `.env.example` and configure the following:

```env
NOTION_API_KEY=your_notion_api_key_here
NOTION_VERSION=2022-06-28
NOTION_API_URL=https://api.notion.com/v1
PORT=8080
```

### Getting Notion API Key

1. Go to [Notion Developers](https://www.notion.so/my-integrations)
2. Create a new integration
3. Copy the API key
4. Share your database with the integration

## Usage

### Running the Server

```bash
go run main.go
```

The server will start on port 8080 (or the port specified in environment variables).

### API Endpoints

#### 1. Search Test Cases
```bash
GET /api/test-cases
```

Returns all test cases that start with `TC_` prefix from your Notion database.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "test_case_key": "01001",
      "page_id": "2946097f-99e0-8057-85ca-f10c7b8d4e68",
      "title": "TC_01001 Login to CMS system by user role in case successfully.",
      "status": "Not started",
      "test_date": "2025-10-22",
      "url": "https://www.notion.so/...",
      "last_edited": "2025-10-22T04:40:00.000Z"
    }
  ],
  "message": "Test cases retrieved successfully"
}
```

#### 2. Get Test Case Blocks
```bash
GET /api/test-cases/{testCaseKey}/blocks
```

Returns all blocks from a specific test case page.

**Parameters:**
- `testCaseKey`: The test case key (e.g., "01001" for "TC_01001")

**Query Parameters:**
- `type=table`: Filter to return only table blocks

**Example:**
```bash
GET /api/test-cases/01001/blocks?type=table
```

**Response:**
```json
{
  "success": true,
  "data": {
    "test_case": {
      "test_case_key": "01001",
      "page_id": "2946097f-99e0-8057-85ca-f10c7b8d4e68",
      "title": "TC_01001 Login to CMS system...",
      "status": "Not started",
      "test_date": "2025-10-22",
      "url": "https://www.notion.so/...",
      "last_edited": "2025-10-22T04:40:00.000Z"
    },
    "blocks": [
      {
        "block_id": "2946097f-99e0-8040-9ed3-c80d828bae02",
        "type": "table",
        "has_children": true,
        "table_info": {
          "table_width": 6,
          "has_column_header": true,
          "has_row_header": false
        }
      }
    ]
  },
  "message": "Blocks retrieved successfully"
}
```

#### 3. Get Block Details
```bash
GET /api/blocks/{blockId}
```

Returns detailed information about a specific block.

**Parameters:**
- `blockId`: The block ID

**Example:**
```bash
GET /api/blocks/2946097f-99e0-8040-9ed3-c80d828bae02
```

## Example Notion Search Query

The application performs the following search against Notion API:

```bash
curl -X POST 'https://api.notion.com/v1/search' \
  -H 'Authorization: Bearer '"$NOTION_API_KEY"'' \
  -H 'Content-Type: application/json' \
  -H 'Notion-Version: 2022-06-28' \
  --data '{
    "query":"External tasks",
    "filter": {
        "value": "data_source",
        "property": "object"
    },
    "sort":{
      "direction":"ascending",
      "timestamp":"last_edited_time"
    }
  }'
```

## Project Structure

```
demo-notion-api/
├── main.go              # Application entry point
├── config/
│   └── config.go        # Configuration management
├── handlers/
│   └── notion.go        # HTTP request handlers
├── services/
│   └── notion.go        # Business logic and Notion API integration
├── models/
│   └── notion.go        # Data structures and models
├── .env.example         # Environment variables template
├── go.mod              # Go module definition
└── README.md           # Project documentation
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `200`: Success
- `400`: Bad Request (missing parameters)
- `404`: Not Found (test case/block not found)
- `500`: Internal Server Error

Error response format:
```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

## Development

### Adding New Endpoints

1. Define new models in `models/notion.go`
2. Implement business logic in `services/notion.go`
3. Create HTTP handlers in `handlers/notion.go`
4. Add routes in `main.go`

### Testing

Test the endpoints using curl or any API client:

```bash
# Test search endpoint
curl http://localhost:8080/api/test-cases

# Test blocks endpoint
curl http://localhost:8080/api/test-cases/01001/blocks?type=table

# Test block details
curl http://localhost:8080/api/blocks/{block-id}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.