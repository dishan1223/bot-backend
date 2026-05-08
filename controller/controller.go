package controller

import (
    "log"
    "os"
    "strings"
    "io"
    "encoding/json"
    "bytes"
    "fmt"
    "net/http"
    "github.com/joho/godotenv"
    "time"
    "github.com/dishan1223/bot-backend/types"
    "github.com/dishan1223/bot-backend/consts"
)



// msg: the message the users want to send to AI Model
// ctx: the context of their previous conversation

// TODO: use a type struct to only store the reply from ai model.
// all other data should not be sent to client
func SendReqToOpenRouter(msg string, ctx string, history []types.Message) ([]types.Message,types.AiResp, error){
    url := "https://openrouter.ai/api/v1/chat/completions"
    err := godotenv.Load()
    if err != nil{
        log.Fatal(err)
    }

    apiKey := os.Getenv("OPENROUTER_API_KEY")

 
    // some extra context based on users prompt. 
    // weather report from OpenWeather
    // time & date data from system. 
    // more will be added
    var system_info string
    if strings.Contains(strings.ToLower(msg), "time"){
        system_info = time.Now().Format(time.RFC3339)
    }
    if strings.Contains(strings.ToLower(msg), "weather")|| 
    strings.Contains(strings.ToLower(msg), "temperature")||
    strings.Contains(strings.ToLower(msg), "temp"){

        w,err := WeatherApiCall()

        if err != nil{
            system_info = ""
        }
        weatherDataforAiModel := fmt.Sprintf("temp: %f, feels_like: %f, humidity: %d, wind_speed: %f ", w.Temp, w.FeelsLike, w.Humidity, w.WindSpeed)

        system_info = weatherDataforAiModel
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




func WeatherApiCall() (types.WeatherReport, error){

    err := godotenv.Load()
    if err != nil{
        log.Fatal(err)
    }


    apiKey := os.Getenv("OPENWEATHER_API_KEY")

    var lat string = "24.4577"
    var lon string = "89.7080"

    var url string = "https://api.openweathermap.org/data/2.5/weather?lat="+lat+"&lon="+lon+"&appid="+apiKey

    req, err := http.NewRequest("GET",url,nil)
    if err != nil{
        log.Fatal(err)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{
        log.Fatal(err)
    }
    defer resp.Body.Close()

    var raw struct{
        Main struct{
            Temp float64 `json:"temp"`
            FeelsLike float64 `json:"feels_like"`
            Humidity int `json:"humidity"`
        } `json:"main"`
        Wind struct{
            Speed float64 `json:"speed"`
        } `json:"wind"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil{
        log.Fatal(err)
    }

    WeatherData := types.WeatherReport{
        Temp: raw.Main.Temp - 273.15,
        FeelsLike: raw.Main.FeelsLike - 273.15,
        Humidity: raw.Main.Humidity,
        WindSpeed: raw.Wind.Speed,
    }


    // weatherData is the actual data and nil is error set to nil
    return WeatherData, nil 
}
