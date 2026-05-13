package types

type Prompt struct{
    Message string `json:"msg"`
    BotId string `json:"botId"`
    UserName string `json:"userName"`
    Lat string `json:"lat"`
    Lon string `json:"lon"`
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


type SearchResult struct {
	Data struct {
		WebPages struct {
			Value []struct {
				Summary string `json:"summary"`
			} `json:"value"`
		} `json:"webPages"`
	} `json:"data"`
}
