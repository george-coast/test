package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	APIKey              string `json:"api_key"`
	BaseURL             string `json:"base_url"`
	OrgID               string `json:"org_id"`
	PasswordManagerEvents struct {
		Enabled    bool     `json:"enabled"`
		EventTypes []string `json:"event_types"`
		LogLevel   string   `json:"log_level"`
	} `json:"password_manager_events"`
}

// Struct to parse the response from JumpCloud API
type PasswordManagerEvent struct {
	ID          string        `json:"id"`
	UUID        string        `json:"uuid"`
	EventType   string        `json:"event_type"`
	Timestamp   string        `json:"timestamp"`
	Service     string        `json:"service"`
	Organization string       `json:"organization"`
	Changes     []Change      `json:"changes"`
	InitiatedBy Initiator     `json:"initiated_by"`
	ClientIP    string        `json:"client_ip"`
	UserAgent   UserAgent     `json:"user_agent"`
	Success     bool          `json:"success"`
	GeoIP       GeoIP         `json:"geoip"`
	ErrorMessage string       `json:"error_message"`
}

type Change struct {
	Field string `json:"field"`
	From  string `json:"from"`
	To    string `json:"to"`
}

type Initiator struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Email string `json:"email"`
}

type UserAgent struct {
	OSName   string `json:"os_name"`
	OSMajor  string `json:"os_major"`
	OSMinor  string `json:"os_minor"`
	Name     string `json:"name"`
}

type GeoIP struct {
	CountryCode  string  `json:"country_code"`
	RegionName   string  `json:"region_name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

func main() {
	// Load the configuration
	config, err := loadConfig("/opt/jumpcloud/config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Check if password manager event logging is enabled
	if config.PasswordManagerEvents.Enabled {
		fmt.Println("Password manager events enabled. Fetching logs...")
		fetchPasswordManagerEvents(config)
	} else {
		fmt.Println("Password manager events are not enabled.")
	}
}

func loadConfig(filePath string) (Config, error) {
	var config Config
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func fetchPasswordManagerEvents(config Config) {
	// Build the API request to fetch password manager logs based on event types
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.BaseURL+"/api/systemlogs/v1/events", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("x-api-key", config.APIKey)
	if config.OrgID != "" {
		req.Header.Add("x-org-id", config.OrgID)
	}

	// Append query parameters for password manager-related event types
	query := req.URL.Query()
	for _, eventType := range config.PasswordManagerEvents.EventTypes {
		query.Add("eventType", eventType)  // Add password manager event type, e.g., "passwordmanager_enable"
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the response body
	var events []PasswordManagerEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Process each password manager event
	for _, event := range events {
		logPasswordManagerEvent(event)
	}
}

// Log the password manager event
func logPasswordManagerEvent(event PasswordManagerEvent) {
	fmt.Printf("Event ID: %s\n", event.ID)
	fmt.Printf("Event Type: %s\n", event.EventType)
	fmt.Printf("Timestamp: %s\n", event.Timestamp)
	fmt.Printf("Service: %s\n", event.Service)
	fmt.Printf("Organization: %s\n", event.Organization)

	fmt.Println("Changes:")
	for _, change := range event.Changes {
		fmt.Printf("\tField: %s, From: %s, To: %s\n", change.Field, change.From, change.To)
	}

	fmt.Printf("Initiated by: %s (%s)\n", event.InitiatedBy.Email, event.InitiatedBy.Type)
	fmt.Printf("Client IP: %s\n", event.ClientIP)
	fmt.Printf("User Agent: %s on OS %s\n", event.UserAgent.Name, event.UserAgent.OSName)
	fmt.Printf("Success: %t\n", event.Success)

	if event.ErrorMessage != "" {
		fmt.Printf("Error: %s\n", event.ErrorMessage)
	}

	fmt.Println("GeoIP Info:")
	fmt.Printf("\tCountry: %s, Region: %s, Lat: %f, Lon: %f\n", event.GeoIP.CountryCode, event.GeoIP.RegionName, event.GeoIP.Latitude, event.GeoIP.Longitude)
	fmt.Println("----------------------")
}
