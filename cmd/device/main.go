package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// parseInt converts a string to an int, returns 0 if invalid
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

var (
	CLIENT_ID     = os.Getenv("GITHUB_CLIENT_ID")
	CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
)

func main() {
	deviceCodeResp, err := startDeviceFlow()
	if err != nil {
		log.Fatalf("❌ Failed to start device flow: %v", err)
	}

	if deviceCodeResp.VerificationURI == "" || deviceCodeResp.UserCode == "" {
		log.Fatal("❌ Invalid response from GitHub: missing verification URI or user code")
	}

	fmt.Println("Visit the following URL and enter the code:")
	fmt.Printf("🔗 %s\n🔑 %s\n\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)

	accessToken, err := pollForToken(deviceCodeResp.DeviceCode, deviceCodeResp.Interval)
	if err != nil {
		log.Fatalf("❌ Failed to retrieve access token: %v", err)
	}
	fmt.Println("✅ Access Token:", accessToken)
}

// Step 1: Start the device flow
func startDeviceFlow() (*DeviceCodeResponse, error) {
	data := fmt.Sprintf("client_id=%s&scope=repo%%20user", CLIENT_ID)
	req, _ := http.NewRequest("POST", "https://github.com/login/device/code", bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub returned status %d: %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &DeviceCodeResponse{
		DeviceCode:      values.Get("device_code"),
		UserCode:        values.Get("user_code"),
		VerificationURI: values.Get("verification_uri"),
		ExpiresIn:       parseInt(values.Get("expires_in")),
		Interval:        parseInt(values.Get("interval")),
	}, nil
}

// Step 2: Poll for the access token
func pollForToken(deviceCode string, interval int) (string, error) {
	fmt.Print("⏳ Polling for the access token")
	defer fmt.Println()
	for {
		fmt.Print(".")
		time.Sleep(time.Duration(interval) * time.Second)

		data := fmt.Sprintf("client_id=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code", CLIENT_ID, deviceCode)
		req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("HTTP error: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var tokenResp TokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			return "", fmt.Errorf("JSON error: %v", err)
		}

		switch tokenResp.Error {
		case "":
			if tokenResp.AccessToken != "" {
				return tokenResp.AccessToken, nil
			}
		case "authorization_pending":
			continue // keep polling
		case "slow_down":
			interval += 5
		default:
			// hard error
			return "", fmt.Errorf("Error: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
		}
	}
}

// Structs

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
