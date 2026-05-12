package systemInfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dishan1223/bot-backend/consts"
	"github.com/dishan1223/bot-backend/internal/api"
	"github.com/dishan1223/bot-backend/internal/service"
)

// trying to make the intent detection system better with this function. this way
// we dont need to write a big condition in each if statements
func containsAny(msg string, keywords []string) bool{
    for _ , k := range keywords {
        if strings.Contains(msg, k){
            return true
        }
    }
    return false
}



// this functions provides extra dependencies to the AI model
// 1. Weather data
// 2. Time & Date data
// this shall be passed asa a context to the AI
// basically a library of data for the ai model to read to provide better response
// maybe on kind of training
func GetModelDeps(msg string) (string, error){

    msg = strings.ToLower(msg)

    var system_info string
    
    // Time&Date intent check
    if containsAny(msg, consts.TIME_KEYWORDS){
        system_info = time.Now().Format(time.RFC3339)
    }
    

    // Weather Intent Check
    if containsAny(msg, consts.WEATHER_KEYWORDS){

        w,err := api.GetWeatherReport()

        if err != nil{
            fmt.Println("Error getting weather data: ",err)
            system_info = ""
            // loggin the error in this function. so maybe i dont need to send the error to 
            // the controller
            return system_info, err 
        }
        weatherDataforAiModel := fmt.Sprintf(
        	`Weather Report:
            Location: %s, %s
            Condition: %s
            Temperature: %.1f°C
            Feels Like: %.1f°C
            Humidity: %d%%
            Wind Speed: %.1f m/s`,
	        w.Location,
	        w.Country,
	        w.Description,
	        w.Temp,
	        w.FeelsLike,
	        w.Humidity,
	        w.WindSpeed,
        )
        system_info +="\nWeather Data Context " + weatherDataforAiModel
    }


    if containsAny(msg, consts.SEARCH_KEYWORDS){
        searchResults, err := service.GetWebResultsFromLangSearch(msg)
        if err != nil{
            fmt.Println("Error from search service")

            return system_info, err
        }

        searchDataToAI , err := sonic.MarshalString(searchResults) 
        if err != nil{
            fmt.Println("Error marshaling search results:", err)
            return system_info, err
        }

        system_info += "\nSearch Context : " + searchDataToAI
    }

    return system_info,nil

}


