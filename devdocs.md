# Developer Documentation - Bot Backend (Nimo)

This document provides a technical overview of the `bot-backend` project for developers.

## Architecture Overview

The backend is built using the [Fiber](https://gofiber.io/) web framework in Go. It follows a modular structure to separate concerns like API communication, internal logic, and type definitions.

### Tech Stack

- **Language:** Go
- **Framework:** Fiber v3
- **AI Integration:** OpenRouter (GPT models)
- **External APIs:** OpenWeather
- **Environment Management:** `godotenv`

## Project Structure

```text
/
├── main.go               # Entry point, route definitions, and middleware
├── go.mod                # Go module dependencies
├── consts/
│   └── consts.go         # Global constants and AI persona definitions
├── controller/
│   └── controller.go     # Business logic for AI interactions
├── internal/
│   ├── api/
│   │   └── getWeatherReport.go # OpenWeather API integration
│   └── systemInfo/
│       └── systemInfo.go # Dynamic context generation (time, weather)
├── types/
│   └── types.go          # Shared data structures and JSON types
└── index.html            # Static landing page (optional)
```

## Setup & Installation

1.  **Clone the repository.**
2.  **Environment Variables:** Create a `.env` file in the root directory with the following keys:
    ```env
    OPENROUTER_API_KEY=your_openrouter_key
    OPENWEATHER_API_KEY=your_openweather_key
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
2.  `main.go` binds the JSON to a `Prompt` struct.
3.  `systemInfo.GetModelDeps` is called to check if the user is asking for time or weather. If so, it fetches the relevant data.
4.  `controller.SendReqToOpenRouter` is called with the user message, persona context, history, and system info.
5.  OpenRouter returns the AI's response.
6.  The conversation history is updated (currently in-memory).
7.  The response is returned to the client.

### 2. Persona Definition

Nimo's persona is defined in `consts/consts.go` as `GeneralAiContext`. It instructs the AI on how to behave, how to address the owner, and what tone to use.

### 3. Dynamic Context Injection

In `internal/systemInfo/systemInfo.go`, the code scans the user's message for keywords like "time", "weather", or "temp". If found, it injects the current server time or real-time weather data from OpenWeather into the system prompt.

## Future Roadmap

-   **Database Migration:** Move from in-memory conversation history to **SQLite** (specifically `sqlite-vec` for vector search capabilities).
-   **Bot Identification:** Better handling of `botId` to segregate conversation history per device.
-   **Dynamic Location:** Allow clients to send latitude and longitude for more accurate weather reporting (currently hardcoded to a specific location in Bangladesh).

## Contributing

1.  Maintain the modular structure.
2.  Ensure new types are added to `types/types.go`.
3.  Add proper error handling and logging for external API calls.
