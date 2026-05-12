package service

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
    "github.com/dishan1223/bot-backend/types"
)
var apiKey string

func InitSearchAPI(apikey string){
    apiKey = apikey
}



// this fuction uses LangSearch api for Web Searching
// this allows our model to have Web Search Capability
func GetWebResultsFromLangSearch(query string) (types.SearchResult, error){
    var LangSearchEndpoint string = "https://api.langsearch.com/v1/web-search"
    var results types.SearchResult

    payload := map[string]any{
        "query": query,
        "summary": true, 
        "count": 3,
    }

    jsonPayload, err := sonic.Marshal(payload)
    if err != nil{
        fmt.Println(err) 
        // change the return format after function is complete
        return results, err 

    }

    req, err := http.NewRequest(http.MethodPost, LangSearchEndpoint, bytes.NewBuffer(jsonPayload))
    if err != nil{
        fmt.Println("Error in request generation : search.go")
        return results, err 
    }

    authorizationHeader := "Bearer " + apiKey
    req.Header.Set("Authorization", authorizationHeader)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil{
        fmt.Println()
        return results, err 
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return results, fmt.Errorf("LangSearch responded with : %d", resp.StatusCode)
    }

    if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&results); err != nil{
        return results, fmt.Errorf("Decode error : %w", err)
    }


    return results, nil
}
