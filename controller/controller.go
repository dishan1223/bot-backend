package controller

import (
    "log"
    "os"
    "io"
    "encoding/json"
    "bytes"
    "net/http"
    "github.com/joho/godotenv"
    "github.com/dishan1223/bot-backend/types"
    "github.com/dishan1223/bot-backend/consts"
    "github.com/dishan1223/bot-backend/internal/systemInfo"
)



// msg: the message the users want to send to AI Model
func SendReqToOpenRouter(msg string, ctx string, history []types.Message) ([]types.Message,types.AiResp, error){
    url := "https://openrouter.ai/api/v1/chat/completions"
    err := godotenv.Load()
    if err != nil{
        log.Fatal(err)
    }

    apiKey := os.Getenv("OPENROUTER_API_KEY")

    system_info, err := systemInfo.GetModelDeps(msg) 
    if err != nil{
        log.Fatal(err)
        system_info = ""
    }
    
    var fullMessages []types.Message

    currentMessage := types.Message{
        Role: "user",
        Content: msg,
    }

    fullMessages = append(fullMessages, types.Message{
        Role: "system",
        Content: consts.GeneralAiContext + " Info: "+system_info,
    })

    fullMessages = append(fullMessages, history...)

    fullMessages = append(fullMessages, currentMessage)

    // "model" here takes the name of the AI model that will be used
    // to generate the response
    payload := map[string]any{
        "model": "openai/gpt-oss-120b:free",
        "messages": fullMessages, 
    }


    jsonData, err := json.Marshal(payload)
    if err != nil{
        log.Fatal(err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil{
        log.Fatal(err)
    }

    var authorizationHeader string = "Bearer " + apiKey
    // set headers
    req.Header.Set("Authorization", authorizationHeader)
    req.Header.Set("content-type", "application/json")


    // send api request to openrouter 
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result types.AiResp
    if err := json.Unmarshal(body, &result); err != nil{
        return nil, result, err
    }

    // add users prompt to prompt history
    newHistory := append(history, currentMessage)
    if len(result.Choices)>0{
        newHistory = append(newHistory, result.Choices[0].Message)
    }

    return newHistory,result, nil 
}





