package hetzner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
	token  string
}

func NewHetznerClient(token string) *Client {
	return &Client{
		client: &http.Client{},
		token:  token,
	}
}

func (c *Client) GetZones() (DnsZones, error) {
	// Get Zones (GET https://dns.hetzner.com/api/v1/zones)

	// Create request

	req, err := http.NewRequest("GET", "https://dns.hetzner.com/api/v1/zones", nil)
	if err != nil {
		return DnsZones{}, err
	}

	// Headers
	req.Header.Add("Auth-API-Token", c.token)

	// Fetch Request
	resp, err := c.client.Do(req)

	if err != nil {
		return DnsZones{}, err
	}

	if resp.StatusCode != 200 {
		return DnsZones{}, err
	}

	zones := DnsZones{}

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &zones)
	if err != nil {
		return DnsZones{}, err
	}

	return zones, nil
}

func (c *Client) GetRecords(zoneID string) (ZoneRecords, error) {
	// Get Records (GET https://dns.hetzner.com/api/v1/records?zone_id={ZoneID})

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://dns.hetzner.com/api/v1/records?zone_id=%s", zoneID), nil)

	// Headers
	req.Header.Add("Auth-API-Token", c.token)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return ZoneRecords{}, err
	}

	// Fetch Request
	resp, err := c.client.Do(req)

	if err != nil {
		return ZoneRecords{}, err
	}

	if resp.StatusCode != 200 {
		return ZoneRecords{}, err
	}

	records := ZoneRecords{}

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &records)
	if err != nil {
		return ZoneRecords{}, err
	}

	return records, nil
}

func (c *Client) UpdateRecord(recordId string, payload []byte) error {
	// Update Record (PUT https://dns.hetzner.com/api/v1/records/{RecordID})

	body := bytes.NewBuffer(payload)

	// Create request
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", recordId), body)

	if err != nil {
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-API-Token", c.token)

	// Fetch Request
	_, err = c.client.Do(req)

	return err
}

func (c *Client) CreateRecord(payload []byte) error {
	// Create Record (POST https://dns.hetzner.com/api/v1/records)

	body := bytes.NewBuffer(payload)

	// Create request
	req, err := http.NewRequest("POST", "https://dns.hetzner.com/api/v1/records", body)

	if err != nil {
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-API-Token", c.token)

	// Fetch Request
	_, err = c.client.Do(req)

	return err
}
