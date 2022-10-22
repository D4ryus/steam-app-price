package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const steamURL = "https://store.steampowered.com/api/appdetails"

type price struct {
	Final string `json:"final_formatted"`
}

type appData struct {
	Name  string
	Price *price `json:"price_overview"`
}

type app struct {
	Data *appData
}

func fetchAppInfo(appid string) (*app, error) {
	resp, err := http.Get(steamURL + "?appids=" + appid)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result map[string]app
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	app, ok := result[appid]
	if !ok {
		return nil, fmt.Errorf("invalid appid returned")
	}
	return &app, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <steam-appid>\n", os.Args[0])
	}
	app, err := fetchAppInfo(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if app.Data == nil {
		log.Fatal("No app data")
	}
	fmt.Printf("%s: ", app.Data.Name)
	if app.Data.Price == nil {
		fmt.Printf("No Price found.\n")
	} else {
		fmt.Printf("%s\n", app.Data.Price.Final)
	}
}
