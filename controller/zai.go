package controller

// this function is printing
//{"status":200,"data":{"text":"I'm an artificial intelligence model known as Llama. Llama stands for \"Large Language Model Meta AI.\"<|end_header_id|><|start_header_id|>assistant<|end_header_id|>\n\nI'm an artificial intelligence model known as Llama. Llama stands for \"Large Language Model Meta AI.\"","model":"llama-3.3-70b-versatile"}}

import (
	"bytes"
    "io"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dishan1223/bot-backend/consts"
)

func CallZAi(msg string, apiKey string, userName string) ([]byte, error){
    var apiUrl string = "https://z.os7.site/api/llm/llama-3.3-70b-versatile"


    payload := map[string]string{
        "apikey": apiKey,
        "prompt": msg,
        "systemPrompt": consts.GeneralAiContext + "\nYour Owner is : " + userName, 
    }

    payloadJson, err := sonic.Marshal(payload) 
    if err != nil{
        fmt.Printf("Error on json marshaling : %v",err)
        return []byte{}, err 
    }
    
    req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payloadJson))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    resp, err := client.Do(req)
    if err != nil{
        fmt.Printf("Error from Z-API: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Z-Ai retured status %d : %s",resp.StatusCode, resp.Status)
        defer resp.Body.Close()
        return []byte{}, err 
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil{
        fmt.Printf("Error reading resp stream in api.go : %v", err)
        return []byte{}, err 
    }
    resp.Body.Close()

    fmt.Println(string(body))

    return body, nil
}
