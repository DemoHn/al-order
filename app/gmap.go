package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// gmap.go: a wrapper for Google Map route API to get distances
// between two points
const gmapAPIRoot = "https://maps.googleapis.com/maps/api/directions/json"

// GMapDirectionResp - standard response of Google Direction API
type GMapDirectionResp struct {
	Status       string
	ErrorMessage string `json:"error_message"`
	Routes       []struct {
		Legs []struct {
			Distance struct {
				Text  string
				Value int
			}
		}
	}
}

// GetRouteDistance - get actual distance from
func GetRouteDistance(loc LocationInfo) (bool, float32, error) {
	apiKey := os.Getenv("GOOGLE_MAP_APIKEY")
	// construct API
	finalURL, _ := url.Parse(gmapAPIRoot)
	startLoc, endLoc := formatPoints(loc)
	params := url.Values{
		"origin":      []string{startLoc},
		"destination": []string{endLoc},
		"key":         []string{apiKey},
	}
	finalURL.RawQuery = params.Encode()

	resp, err := fetchData(finalURL.String())
	if err != nil {
		return false, 0.0, err
	}

	// get json resp
	var jsonResp GMapDirectionResp
	if err := json.Unmarshal(resp, &jsonResp); err != nil {
		return false, 0.0, fmt.Errorf("parse json error: %s", err.Error())
	}

	if jsonResp.Status == "OK" {
		distance := jsonResp.Routes[0].Legs[0].Distance.Value
		return true, float32(distance), nil
	}

	// if
	if jsonResp.Status == "ZERO_RESULTS" {
		return false, 0.0, nil
	}

	return false, 0.0, fmt.Errorf("unknownError from API response:%s, %s", jsonResp.Status, jsonResp.ErrorMessage)
}

// private function

// formatPoints
// example format: -12.25,10.25
// return (start, end)
func formatPoints(loc LocationInfo) (string, string) {
	startLoc := fmt.Sprintf("%s,%s", loc.StartLat, loc.StartLng)
	endLoc := fmt.Sprintf("%s,%s", loc.EndLat, loc.EndLng)

	return startLoc, endLoc
}

func fetchData(finalURL string) ([]byte, error) {
	resp, err := http.Get(finalURL)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// fetch data
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return respBody, nil
}
