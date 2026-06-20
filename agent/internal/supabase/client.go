package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	URL    string
	Key    string
	UserID string
}

type ActivityLog struct {
	UserID     string     `json:"user_id"`
	DeviceID   string     `json:"device_id"`
	AppName    string     `json:"app_name"`
	Title      string     `json:"title"`
	URL        string     `json:"url,omitempty"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}

func NewClient(url, key string) *Client {
	return &Client{
		URL: url,
		Key: key,
	}
}

func (c *Client) InsertLog(log ActivityLog) error {
	endpoint := fmt.Sprintf("%s/rest/v1/activity_logs", c.URL)
	
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.Key)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("supabase request failed with status: %s", resp.Status)
	}

	return nil
}
