package systemInfo

import(
    "strings"
    "time"
    "fmt"
    "github.com/dishan1223/bot-backend/internal/api" 
)


// this functions provides extra dependencies to the AI model
// 1. Weather data
// 2. Time & Date data
// this shall be passed asa a context to the AI
// basically a library of data for the ai model to read to provide better response
// maybe on kind of training
func GetModelDeps(msg string) (string, error){
    var system_info string
    if strings.Contains(strings.ToLower(msg), "time"){
        system_info = time.Now().Format(time.RFC3339)
    }
    if strings.Contains(strings.ToLower(msg), "weather")|| 
    strings.Contains(strings.ToLower(msg), "temperature")||
    strings.Contains(strings.ToLower(msg), "temp"){
        w,err := api.GetWeatherReport()

        if err != nil{
            fmt.Println("Error getting weather data: ",err)
            system_info = ""
        }
        weatherDataforAiModel := fmt.Sprintf("temp: %f, feels_like: %f, humidity: %d, wind_speed: %f ", w.Temp, w.FeelsLike, w.Humidity, w.WindSpeed)

        system_info = weatherDataforAiModel
    }

    return system_info,nil

}


