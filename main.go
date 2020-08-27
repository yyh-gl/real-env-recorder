package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var apiKey = os.Getenv("REMO_API_KEY")

type realEnvInfo struct {
	Temperature float64
	Humidity    float64
	Illuminance float64
	Motion      float64
}

// fetchRealEnvInfo gets real environment info from Nature Remo API
func fetchRealEnvInfo() (*realEnvInfo, error) {
	type response []struct {
		ID                string  `json:"id"`
		Name              string  `json:"name"`
		TemperatureOffset float64 `json:"temperature_offset"`
		HumidityOffset    float64 `json:"humidity_offset"`
		FirmwareVersion   string  `json:"firmware_version"`
		MacAddress        string  `json:"mac_address"`
		SerialNumber      string  `json:"serial_number"`
		NewestEvents      struct {
			Temperature struct {
				Val       float64   `json:"val"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"te"`
			Humidity struct {
				Val       float64   `json:"val"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"hu"`
			Illuminance struct {
				Val       float64   `json:"val"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"il"`
			Motion struct {
				Val       float64   `json:"val"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"mo"`
		} `json:"newest_events"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.nature.global/1/devices", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r response
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &realEnvInfo{
		Temperature: r[0].NewestEvents.Temperature.Val + r[0].TemperatureOffset,
		Humidity:    r[0].NewestEvents.Humidity.Val + r[0].HumidityOffset,
		Illuminance: r[0].NewestEvents.Illuminance.Val,
		Motion:      r[0].NewestEvents.Motion.Val,
	}, nil
}

func main() {
	info, err := fetchRealEnvInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("========================")
	fmt.Println(info)
	fmt.Println("========================")
}
