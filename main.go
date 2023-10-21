package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempF     float64 `json:"temp_f"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	q := "Salt Lake City"
	api := os.Getenv("KEY")
	if api == "" {
		log.Fatal("API key not found in environment variables")
	}
	fmt.Printf("this is %T\n", os.Getenv("KEY"))
	_, lku := os.LookupEnv("KEY")
	fmt.Println("apikey exist?:", lku)

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	resp, err := http.Get(
		"http://api.weatherapi.com/v1/forecast.json?key=" + api + "&q=" + q + "&days=1&aqi=no&alerts=no",
	)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != 200 {
		log.Fatal("API unavailable")
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var weather Weather
	err = json.Unmarshal(output, &weather)
	if err != nil {
		log.Fatal(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempF,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempF,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}
