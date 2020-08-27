package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type realEnvInfo struct {
	Temperature string
	Humidity    string
	Illuminance string
	Motion      string
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
	req.Header.Set("Authorization", "Bearer "+os.Getenv("REMO_API_KEY"))

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
		Temperature: strconv.FormatFloat(r[0].NewestEvents.Temperature.Val+r[0].TemperatureOffset, 'f', 2, 64),
		Humidity:    strconv.FormatFloat(r[0].NewestEvents.Humidity.Val+r[0].HumidityOffset, 'f', 2, 64),
		Illuminance: strconv.FormatFloat(r[0].NewestEvents.Illuminance.Val, 'f', 2, 64),
		Motion:      strconv.FormatFloat(r[0].NewestEvents.Motion.Val, 'f', 2, 64),
	}, nil
}

// writeToCSV writes real environment info to CSV
func writeToCSV(info realEnvInfo) error {
	f, err := os.OpenFile(os.Getenv("CSV_FILE"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	w := csv.NewWriter(f)
	if err := w.Write([]string{time.Now().Format(time.RFC3339), info.Temperature, info.Humidity, info.Illuminance, info.Motion}); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func main() {
	info, err := fetchRealEnvInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := writeToCSV(*info); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("===== RESULT =====")
	fmt.Println("気温:", info.Temperature)
	fmt.Println("湿度:", info.Humidity)
	fmt.Println("照度:", info.Illuminance)
	fmt.Println("人感:", info.Motion)
	fmt.Println("===== RESULT =====")
}
