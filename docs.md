# Nimo AI Backend Documentation

Welcome to the documentation for the Nimo AI Backend. Nimo is a digital AI pet companion designed to be cute, helpful, and interactive.

## Overview

The Nimo Backend provides a RESTful API for communicating with the Nimo AI. It handles message processing, persona management, and integrates real-time data like weather, time, and web search results to provide contextual responses.

## API Reference

### Chat Endpoint

Send a message to Nimo and receive a response.

*   **URL:** `https://nimo-api.onrender.com/api/v1/chat`
*   **Method:** `POST`
*   **Content-Type:** `application/json`

> **Note:** This API is currently hosted on Render's free plan for testing purposes. If the service hasn't been used for a while, it may take 30-60 seconds to "spin up" on the first request.

#### Request Body

| Field | Type | Description |
| :--- | :--- | :--- |
| `msg` | `string` | The message or question for Nimo. |
| `botId` | `string` | A unique identifier for the specific bot/device. |
| `userName` | `string` | The name of the user/owner. |
| `lat` | `string` | Latitude for weather services. |
| `lon` | `string` | Longitude for weather services. |

**Example Request:**

```json
{
    "msg": "Nimo, what is the weather today?",
    "botId": "bot_12345",
    "userName": "Ishtiaq",
    "lat": "24.4577",
    "lon": "89.7080"
}
```

#### Response Body

| Field | Type | Description |
| :--- | :--- | :--- |
| `reply` | `string` | Nimo's response text. |
| `status` | `string` | Status of the request (e.g., "success"). |

**Example Response:**

```json
{
    "reply": "It's currently 25°C and sunny in Sirajganj! Perfect for a walk, isn't it?",
    "status": "success"
}
```

## Features

*   **Contextual Awareness:** Nimo knows the current date and time.
*   **Real-time Weather:** Fetches live weather data based on the coordinates provided in the request.
*   **Web Search:** Can perform web searches via LangSearch to answer questions about recent events or complex topics.
*   **Persona-Driven:** Nimo has a unique personality—cute, helpful, and focused on his owner.
*   **Language Support:** Nimo replies naturally in the same language the user uses.
*   **Study Helper:** Specifically tuned to be a study companion in the context of Bangladesh.

## Intent Detection

Nimo uses simple keyword detection to determine when to fetch real-time information. Note that more sophisticated intent detection will be available on full launch.

### Keywords

| Intent | Trigger Keywords |
| :--- | :--- |
| **Weather** | `weather`, `rain`, `rains`, `cold`, `hot`, `temperature`, `temp` |
| **Web Search** | `search`, `web`, `website`, `news`, `story` |

## Interaction Examples

**1. Asking for time:**
*   User: "Nimo, what time is it?"
*   Nimo: "It's currently 3:45 PM, Ishtiaq!"

**2. Asking for weather:**
*   User: "Will it rain today?"
*   Nimo: "The forecast says it's mostly clear with no rain expected. Enjoy your day!"

**3. Web Search:**
*   User: "Search for the latest news about SpaceX."
*   Nimo: "I found that SpaceX recently launched another batch of Starlink satellites. They are so busy!"

## Error Handling

The API returns standard HTTP status codes:

*   `200 OK`: Request was successful.
*   `400 Bad Request`: Invalid request payload.
*   `500 Internal Server Error`: Something went wrong on the server side.
