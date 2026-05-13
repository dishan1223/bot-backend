package api

import (
    "github.com/bytedance/sonic"
	"fmt"
	"net/http"
    "time"

	"github.com/dishan1223/bot-backend/types"
)

var apiKey string

// lat and lon will be embedded in the esp32 device. uses will be able to change it
// via our Mobile Application and send it to this server
var lat string
var lon string


// @info which this function we dont have to call api key from .env everytime this function runs
func InitWeatherAPI(key string, l string, lo string){
    apiKey = key
    lat = l
    lon = lo
}

// @info this is the provider function that fetches data from openweather and returns
// it to the AI via systemInfoContext (internal/systemInfo/systemInfo.go)
func GetWeatherReport() (types.WeatherReport, error){

    var url string = fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s",
		lat,
		lon,
		apiKey,
	)

    // "GET" -> http.MethodGet. Read the docs yesterday
    req, err := http.NewRequest(http.MethodGet,url,nil)
    if err != nil{
        fmt.Printf("failed to create request to get weather data: %e",err)
        return types.WeatherReport{}, err
    }

    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    resp, err := client.Do(req)
    if err != nil{
        fmt.Printf("failed to send api request to openweather: %e", err)
        return types.WeatherReport{}, err
    }
    defer resp.Body.Close()

    // if api sends a non-200 response then no need to run the rest of the code
    if resp.StatusCode != http.StatusOK{
        fmt.Printf("failed to get weather data: %s", resp.Status)
        return types.WeatherReport{}, nil
    }

    var raw struct {
		Name string `json:"name"`

		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`

		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`

		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`

		Sys struct {
			Country string `json:"country"`
		} `json:"sys"`
	} 

    // Sonic is a blazingly fast JSON serializing & deserializing library, accelerated by JIT (just-in-time compiling)
    // and SIMD (single-instruction-multiple-data). 
    // this is relatively much faster than standard encoding/json library
    SonicDecoder := sonic.ConfigDefault.NewDecoder(resp.Body) 
    if err := SonicDecoder.Decode(&raw); err != nil{
        fmt.Printf("failed to decode opendata api response : %e", err)
        return types.WeatherReport{}, err
    }

    weatherData := types.WeatherReport{
		Temp:       raw.Main.Temp - 273.15,
		FeelsLike:  raw.Main.FeelsLike - 273.15,
		Humidity:   raw.Main.Humidity,
		WindSpeed:  raw.Wind.Speed,
		Location:   raw.Name,
		Country:    raw.Sys.Country,
		Description: func() string {
			if len(raw.Weather) > 0 {
				return raw.Weather[0].Description
			}
			return ""
		}(),
	} 


    // weatherData is the actual data and nil is error set to nil
    return weatherData, nil 
}
