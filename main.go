package main

// NOTE: 
// WE have many error handling right now. To track them in production, 
// we will use a Telegram Bot who'll notify us on every error in the server


import (
	"fmt"
	"os"

	"github.com/dishan1223/bot-backend/consts"
	"github.com/dishan1223/bot-backend/controller"
	"github.com/dishan1223/bot-backend/internal/api"
	"github.com/dishan1223/bot-backend/internal/service"
	"github.com/dishan1223/bot-backend/types"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
)

// Get this history from vector database, based on botId
// botId is a unique identifier for each bot that is embedded in the esp32 device
var history = []types.Message{
    {
        Role: "user",
        Content: "hello how are you Nimo?",
    },
    {
        Role: "assistant",
        Content: "Hey I am fine what about you?",
    },
}


func main(){


    // Loading environment variables
    // and initializing dependencies
    if err := godotenv.Load(); err != nil{
        fmt.Println("Error adding the environment variables")
    }
    apiKey := os.Getenv("OPENROUTER_API_KEY")
    service.InitSearchAPI(os.Getenv("LANGSEARCH_API_KEY"))



    // backend part starts from here
    app := fiber.New()
    v1 := app.Group("/api/v1")
    

    // gofiber middlewares
    app.Use(cors.New())
    app.Use(requestid.New())
    app.Use(logger.New(logger.Config{
        Format: "$[${ip}]:${port} | {pid} ${requestid} ${status} - ${method} ${path}\n",
    }))

    // domain:3000/api/v1/chat
    v1.Post("/chat", func(c fiber.Ctx) error {
        p := new(types.Prompt)
        // initialize weather api service
        api.InitWeatherAPI(os.Getenv("OPENWEATHER_API_KEY"), p.Lat, p.Lon)
       
        // json example from client
        //  {
        //      "msg": "Nimo whats todays date? ",
        //      "botId": "ab123x",
        //      "userName": "Ishtiaq Dishan",
        //      "lat": "24.4577",
        //      "lon": "89.7080"
        // }
        // ctx is saved in consts/consts.go file.(development)
        // No need to take contexts from the esp32. 
        // We will use an SQLite database to store user conversations and contexts

        if err := c.Bind().Body(p); err != nil{
            fmt.Println(err)
        }

        // free up context window.
        // NOTE : This is only for the alpha version of this codebase for testing purposes.
        // In production we will use a database and conditionally take contexts
        if len(history) >= consts.CONTEXT_WINDOW{
            history = history[consts.TRIM_COUNT:]
        }
        

        // send request to OpenRouter. newHistory is the updated conversational history
        // we will need to update our database with this newHistory for each bot.
        // each bot's id is the botId
        newHistory,fullResponse, err := controller.SendReqToOpenRouter(p.Message, consts.GeneralAiContext, history, apiKey, p.UserName)
        if err != nil {
            fmt.Println(err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Internal server error",
            })
        }

        // updating conversational history
        // NOTE: currently we are not using any database. But we will move to sqlite-vec
        // as our primary vector database
        history = newHistory
        replyText := ""
        if len(fullResponse.Choices) > 0 {
            replyText = fullResponse.Choices[0].Message.Content
        }

        c.Set("Content-Type", "application/json")
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "reply": replyText,
            "status": "success",
        })
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    fmt.Println(app.Listen(":" + port))
}
