package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		Temp_c    float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	// if country not provided, use this as default
	country := "ethiopia"
	godotenv.Load()

	weather_api_key := os.Getenv("WEATHER_API_KEY")

	if len(os.Args) > 2 {
		country = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + weather_api_key + "&q=" + country + "&days=7")

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic("api connetcion unsuccessful")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)

	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("-------------------------\n")
	fmt.Printf(
		"%s, %s | %.fC, %s\n",
		location.Name,
		location.Country,
		current.Temp_c,
		current.Condition.Text)
	fmt.Printf("-------------------------\n\n")

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		msg := fmt.Sprintf("%s > %.fC %.f%% %s\n",
			date.Format("02/Jan/2006"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text)

		if hour.ChanceOfRain < 80 {
			fmt.Print(msg)
		} else {
			color.Red(msg)
		}
	}
}
