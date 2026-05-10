package api

import (
	"fmt"
	"net/http"
    "encoding/json"
)

// Note: this api does not work so good. we will be using gnews though its not free.

// this is a free news api that I found on github.
// Build by : Abu Nayim Faisal (https://github.com/faisal-shohag)
// Api github repo: https://github.com/faisal-shohag/news-api
// this API scraps BBC Bengali website

// NOTE:
// This scraper targets https://www.bbc.com/bengali. Ensure compliance with BBC's terms of service and robots.txt

func GetNewsReports() ([]string, error){
    var url string = "https://news-api-fs.vercel.app/api/popular"
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil{
        return nil, err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        return nil, err
    }
    defer resp.Body.Close()
    
    var apiResponse struct {
        Articles []struct {
            Title string `json:"title"`
        } `json:"articles"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil{
        fmt.Printf("failed to decode opendata api response : %e", err)
    }
    var titles []string
    for _, item := range apiResponse.Articles {
        titles = append(titles, item.Title)
    }
    
    

    return titles, nil
}
