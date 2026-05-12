package types

type Prompt struct{
    Message string `json:"msg"`
    BotId string `json:"botId"`
    UserName string `json:"userName"`
}

type WeatherReport struct {
	Temp        float64
	FeelsLike   float64
	Humidity    int
	WindSpeed   float64
	Location    string
	Country     string
	Description string
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

