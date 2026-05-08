package types

type Prompt struct{
    Message string `json:"msg"`
}

type WeatherReport struct{
    Temp float64 `json:"temp"`
    FeelsLike float64 `json:"feels_like"`
    Humidity int `json:"humidity"`
    WindSpeed float64 `json:"wind_speed"`
}

type Message struct{
    Role string `json:"role"`
    Content string `json:"content"`
}


// this type is for getting response
type AiResp struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

