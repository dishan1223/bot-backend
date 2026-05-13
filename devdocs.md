# Developer Documentation - Bot Backend (Nimo)

This document provides a technical overview of the `bot-backend` project for developers.

## Architecture Overview

The backend is built using the [Fiber](https://gofiber.io/) web framework in Go. It follows a modular structure to separate concerns like API communication, internal logic, and type definitions.

### Tech Stack

- **Language:** Go 1.25.0
- **Framework:** Fiber v3 (High performance, minimalist)
- **AI Integration:** OpenRouter (using `openai/gpt-oss-120b:free` model)
- **External APIs:** 
    - **OpenWeather:** Real-time weather data.
    - **LangSearch:** Web searching capabilities.
- **JSON Processing:** `github.com/bytedance/sonic` (SIMD-accelerated, blazingly fast)
- **Environment Management:** `godotenv`

## Project Structure

```text
/
├── main.go               # Entry point, route definitions, and middleware
├── go.mod                # Go module dependencies
├── consts/
│   └── consts.go         # Global constants, keywords, and AI persona
├── controller/
│   └── controller.go     # Business logic for OpenRouter interactions
├── internal/
│   ├── api/
│   │   ├── getWeatherReport.go # OpenWeather API integration
│   │   └── newsReport.go       # BBC Bengali news scraper (Partially integrated)
│   ├── service/
│   │   └── search.go           # LangSearch API integration for web search
│   └── systemInfo/
│       └── systemInfo.go       # Intent detection and dynamic context generation
├── types/
│   └── types.go          # Shared data structures and JSON types
└── index.html            # Static landing page
```

## Setup & Installation

1.  **Clone the repository.**
2.  **Environment Variables:** Create a `.env` file in the root directory with the following keys:
    ```env
    OPENROUTER_API_KEY=your_openrouter_key
    OPENWEATHER_API_KEY=your_openweather_key
    LANGSEARCH_API_KEY=your_langsearch_key
    PORT=3000
    ```
3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```
4.  **Run the Server:**
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:3000`.

## Core Logic

### 1. Request Flow

1.  User sends a POST request to `/api/v1/chat`.
2.  `main.go` binds the JSON to a `Prompt` struct using `c.Bind().Body(p)`.
3.  `api.InitWeatherAPI` is called with coordinates from the request.
4.  `controller.SendReqToOpenRouter` orchestrates the AI interaction:
    - Calls `systemInfo.GetModelDeps` to gather dynamic context.
    - Injects the persona (`GeneralAiContext`), user name, and dynamic context into the system message.
    - Sends the request to OpenRouter with conversation history.
5.  `systemInfo.GetModelDeps` performs **Intent Detection**:
    - Scans message for keywords defined in `consts/consts.go`.
        - **Weather:** `weather`, `rain`, `rains`, `cold`, `hot`, `temperature`, `temp`
        - **Search:** `search`, `web`, `website`, `news`, `story`
    - **Time/Date:** Returns current server time in RFC3339 format.
    - **Weather:** Calls `api.GetWeatherReport` to fetch live data.
    - **Search:** Calls `service.GetWebResultsFromLangSearch` for web search results.
6.  The response is returned and conversation history is updated in-memory.

### 2. High Performance JSON

The project uses the `sonic` library for JSON marshaling and unmarshaling. This provides a significant performance boost over the standard `encoding/json` library, especially useful for handling large AI contexts and responses.

### 3. Context Window Management

Currently, the server maintains a global `history` slice in memory.
- `CONTEXT_WINDOW` (40) defines the maximum number of messages kept.
- `TRIM_COUNT` (10) defines how many messages are removed when the window is exceeded.
- *Note: This is a temporary solution for the alpha version.*

## Future Roadmap

-   **Database Migration:** Move from in-memory conversation history to **SQLite** (specifically `sqlite-vec` for vector search capabilities).
-   **Multi-tenant History:** Segregate conversation history per `botId` using the database.
-   **News Integration:** Fully integrate `newsReport.go` into the `GetModelDeps` workflow.
-   **Error Notification:** Implement a Telegram Bot to notify developers of server errors in real-time.

## Contributing

1.  Maintain the modular structure.
2.  Use `sonic` for all JSON operations.
3.  Add new intent keywords to `consts/consts.go`.
4.  Ensure all new dependencies are documented in `go.mod`.
