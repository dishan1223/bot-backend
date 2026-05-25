# Nimo AI Backend

Nimo is a digital AI pet companion backend designed for contextual interaction and helpfulness. Built with Go using the Fiber v3 framework, this service integrates real-time environmental data including weather, temporal context, and web search capabilities to provide informed responses.

## Core Features

- **Context-Aware Interaction:** Automatically detects user intent to provide real-time information.
- **Environmental Data Integration:**
    - **Weather Services:** Integrated with OpenWeather API for live meteorological updates.
    - **Web Search:** Real-time information retrieval via LangSearch API.
    - **Temporal Awareness:** Native server-side time and date synchronization.
- **Multilingual Support:** Dynamic language adaptation based on user input.
- **Performance Optimized:** Utilizes `bytedance/sonic` for high-speed SIMD-accelerated JSON processing.
- **Modular Architecture:** Clean separation of concerns following Fiber v3 best practices.

## Project Structure

```text
.
├── main.go               # Application entry point and routing
├── consts/
│   └── consts.go         # Global constants and configuration
├── controller/
│   └── controller.go     # Core business logic and AI orchestration
├── internal/
│   ├── api/
│   │   ├── getWeatherReport.go # Weather service integration
│   │   └── newsReport.go       # News aggregation logic
│   ├── initializer/
│   │   └── LoadEnv.go          # Environment configuration loader
│   ├── service/
│   │   └── search.go           # External search service implementation
│   └── systemInfo/
│       └── systemInfo.go       # Intent detection and dynamic context synthesis
├── types/
│   └── types.go          # Shared data structures and schemas
└── index.html            # Application landing page
```

## Documentation

For more detailed information, please refer to the following documentation files:
- `docs.md`: Comprehensive API reference and feature overview.
- `devdocs.md`: Technical architecture and development guide.

## Setup and Installation

### Prerequisites

- Go 1.25.0 or higher
- Required API Keys:
    - OpenRouter (LLM orchestration)
    - OpenWeather (Meteorological data)
    - LangSearch (Web search capabilities)

### Installation Steps

1. **Clone the repository:**
   ```bash
   git clone https://github.com/dishan1223/bot-backend.git
   cd bot-backend
   ```

2. **Environment Configuration:**
   Create a `.env` file in the root directory with the following parameters:
   ```env
   OPENROUTER_API_KEY=your_openrouter_key
   OPENWEATHER_API_KEY=your_openweather_key
   LANGSEARCH_API_KEY=your_langsearch_key
   PORT=3000
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Initialize Server:**
   ```bash
   go run main.go
   ```
   The service will be accessible at `http://localhost:3000`.

## API Documentation

### Chat Interface
- **Endpoint:** `POST /api/v1/chat`
- **Description:** Processes user messages and returns context-aware AI responses.
- **Request Schema:**
  ```json
  {
    "message": "Current weather in London?",
    "userName": "User",
    "lat": "51.5074",
    "lon": "0.1278"
  }
  ```

### Internal Services
- **Endpoint:** `GET /api/v1/zai`
- **Description:** Internal diagnostic endpoint for ZAi service validation.

## License

Internal use only. All rights reserved.
