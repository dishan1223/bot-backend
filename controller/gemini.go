package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/dishan1223/bot-backend/consts"
	"github.com/dishan1223/bot-backend/internal/systemInfo"
	"github.com/dishan1223/bot-backend/types"
	"github.com/gofiber/fiber/v3"
	"google.golang.org/genai"
)

// SendReqToGemini sends a request to the Google Gemini API using the new unified SDK.
func SendReqToGemini(msg string, history []types.Message, apiKey string, userName string, lat string, lon string) ([]types.Message, string, error) {
	ctx := context.Background()

	// Initialize the client pointing to the Gemini Developer API backend
	// The SDK automatically detects GOOGLE_API_KEY from environment, 
	// but we can also set it if needed. For consistency, we'll assume it's in the env.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Gemini client: %v", err)
	}

	system_info, err := systemInfo.GetModelDeps(msg, lat, lon)
	if err != nil {
		fmt.Printf("Error getting systemInfo: %v\n", err)
		system_info = ""
	}

	// Prepare system instruction
	systemContext := consts.GeneralAiContext + "\nYour Owner is : " + userName + "\nMore Contexts: " + system_info
	
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemContext}},
		},
	}

	// Prepare contents (history + current message)
	var contents []*genai.Content

	// Add history
	for _, m := range history {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: m.Content}},
		})
	}

	// Add current message
	contents = append(contents, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: msg}},
	})

	// Using gemini-2.5-flash as explicitly requested by the user.
	modelName := "gemini-2.5-flash"
	
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate content with Gemini: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, "", fmt.Errorf("empty response from Gemini")
	}

	replyText := resp.Text()

	// Prepare new history
	newHistory := append(history, types.Message{
		Role:    "user",
		Content: msg,
	})
	newHistory = append(newHistory, types.Message{
		Role:    "assistant",
		Content: replyText,
	})

	return newHistory, replyText, nil
}

// GeminiHandler is the GoFiber API handler for Gemini chat.
func GeminiHandler(c fiber.Ctx) error {
	p := new(types.Prompt)

	if err := c.Bind().Body(p); err != nil {
		fmt.Println("Binding error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gemini API key not configured",
		})
	}

	// We'll use nil for history for now.
	newHistory, replyText, err := SendReqToGemini(p.Message, nil, apiKey, p.UserName, p.Lat, p.Lon)
	if err != nil {
		fmt.Println("Gemini error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	_ = newHistory

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"reply":  replyText,
		"status": "success",
	})
}
