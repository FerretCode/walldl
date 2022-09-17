package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/akamensky/argparse"
	"github.com/caarlos0/env/v6"
)

type Wallpapers struct {
	Wallpapers []Wallpaper `json:"photos"`
}

type Wallpaper struct {
	ID int `json:"id"`
	URL string `json:"path"`
	Source Source `json:"src"`
}

type Source struct {
	Original string `json:"original"`
}

type ApiKey struct {
	PexelsApiKey string `env:"WALLDL_API_KEY"`
}

func main() {
	parser := argparse.NewParser("walldl", "Downloads wallpapers in a category")

	category := parser.String("c", "category", &argparse.Options{Required: true, Help: "The category of wallpaper"})
	number := parser.Int("n", "number", &argparse.Options{Required: true, Help: "The number of wallpapers to fetch"})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(err)
	}

	var apiKey ApiKey

	if err := env.Parse(&apiKey); err != nil {
		log.Fatal(err)
	}

	getWallpapers(*category, *number, apiKey)
}

func getWallpapers(category string, number int, apiKey ApiKey) {
	if apiKey.PexelsApiKey == "" {
		log.Fatal("Please set the WALLDL_API_KEY environment variable to your Pexels API key!")
	}
	
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET", 
		fmt.Sprintf(
			"https://api.pexels.com/v1/search?query=%s&per_page=%s", 
			url.QueryEscape(category),
			strconv.Itoa(number),
		),
		nil,
	)

	if err != nil {
		fmt.Printf("There was an error fetching the wallpapers: %s", err)

		return
	}

	req.Header.Add("Authorization", apiKey.PexelsApiKey)

	res, err :=	client.Do(req) 

	if err != nil {
		fmt.Printf("There was an error fetching the wallpapers: %s", err)

		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("There was an error fetching the wallpapers: %s", err)
	}

	var wallpapers Wallpapers

	jsonErr := json.Unmarshal(body, &wallpapers)

	if jsonErr != nil {
		fmt.Printf("There was an error fetching the wallpapers: %s", jsonErr)
	

		return
	}

	res.Body.Close()

	for _, v := range wallpapers.Wallpapers {
		image, err := http.Get(v.Source.Original)

		if err != nil {
			log.Fatal(err)
		}

		defer image.Body.Close()

		home, err := os.UserHomeDir()

		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(fmt.Sprintf("%s/walldl/wallpapers/%s", home, strconv.Itoa(v.ID)))

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		_, err = io.Copy(file, image.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Image %s downloaded\n", strconv.Itoa(v.ID))
	}
}
