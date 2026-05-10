package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dishan1223/bot-backend/types"
)

var apiKey string


// @info which this function we dont have to call api key from .env everytime this function runs
func InitWeatherAPI(key string){
    apiKey = key
}

// @info this is the provider function that fetches data from openweather and returns
// it to the AI via systemInfoContext (internal/systemInfo/systemInfo.go)
func GetWeatherReport() (types.WeatherReport, error){

    // lat and lon will be embedded in the esp32 device. uses will be able to change it
    // via our Mobile Application and send it to this server
    var lat string = "24.4577"
    var lon string = "89.7080"

    var url string = "https://api.openweathermap.org/data/2.5/weather?lat="+lat+"&lon="+lon+"&appid="+apiKey

    req, err := http.NewRequest("GET",url,nil)
    if err != nil{
        fmt.Printf("failed to create request to get weather data: %e",err)
        return types.WeatherReport{}, err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{
        fmt.Printf("failed to send api request to openweather: %e", err)
        return types.WeatherReport{}, err
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
        fmt.Printf("failed to decode opendata api response : %e", err)
        return types.WeatherReport{}, err
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
