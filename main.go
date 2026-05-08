package main

import(
    "log"
    "fmt"
    "github.com/gofiber/fiber/v3"
    "github.com/dishan1223/bot-backend/types"
    "github.com/dishan1223/bot-backend/consts"
    "github.com/dishan1223/bot-backend/controller"
)

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

    app := fiber.New()
    v1 := app.Group("/api/v1")

    app.Get("/", func( c fiber.Ctx ) error {
        return c.SendString("hello from backend")
    })

    v1.Get("/weather", func(c fiber.Ctx) error {
        
        data,err := controller.WeatherApiCall()
        if err != nil{
            log.Fatal(err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to fetch weather data",
            })
        }

        c.Set("Content-Type","application/json")
        return c.Status(fiber.StatusOK).JSON(data)
    })

    // domain:3000/api/v1/chat
    v1.Post("/chat", func(c fiber.Ctx) error {
        p := new(types.Prompt)
       
              

        // this ctx shall be embeded to the bot itself
        // json example from client
        //  {
        // "msg": "Nimo whats todays date? ",
        // }
        // ctx is saved in consts/consts.go file.(development)
        // No need to take contexts from the esp32. 
        // We will use an SQLite database to store user conversations and contexts

        if err := c.Bind().Body(p); err != nil{
            fmt.Println(err)
        }

        newHistory,fullResponse, err := controller.SendReqToOpenRouter(p.Message, consts.GeneralAiContext, history)
        if err != nil {
            log.Fatal(err)
            c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Internal server error",
            })
        }

        // updating conversational history
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

    fmt.Println(app.Listen(":3000"))
}
