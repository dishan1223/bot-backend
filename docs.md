# Nimo AI Backend Documentation

Welcome to the documentation for the Nimo AI Backend. Nimo is a digital AI pet companion designed to be cute, helpful, and interactive.

## Overview

The Nimo Backend provides a RESTful API for communicating with the Nimo AI. It handles message processing, persona management, and integrates real-time data like weather and time to provide contextual responses.

## API Reference

### Chat Endpoint

Send a message to Nimo and receive a response.

*   **URL:** `/api/v1/chat`
*   **Method:** `POST`
*   **Content-Type:** `application/json`

#### Request Body

| Field | Type | Description |
| :--- | :--- | :--- |
| `msg` | `string` | The message or question for Nimo. |
| `botId` | `string` | A unique identifier for the specific bot/device. |

**Example Request:**

```json
{
    "msg": "Nimo, what is the weather today?",
    "botId": "bot_12345"
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
    "reply": "It's currently 25°C and sunny! Perfect for a walk, isn't it?",
    "status": "success"
}
```

## Features

*   **Contextual Awareness:** Nimo knows the current date, time, and weather.
*   **Persona-Driven:** Nimo has a unique personality—cute, helpful, and focused on his owner.
*   **Language Support:** Nimo replies in the same language the user uses.
*   **Study Helper:** Specifically tuned to be a study companion in the context of Bangladesh.

## Error Handling

The API returns standard HTTP status codes:

*   `200 OK`: Request was successful.
*   `400 Bad Request`: Invalid request payload.
*   `500 Internal Server Error`: Something went wrong on the server side.
