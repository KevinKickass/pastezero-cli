package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"log"
)

type handshakeResponse struct {
	ClientID  string `json:"client_id"`
	Signature string `json:"signature"`
}

func EnsureClientID(apiURL string) (*Config, error) {
	// Config bereits vorhanden?
	cfg, err := LoadConfig()
	if err == nil {
		return cfg, nil
	}

	// Handshake
	res, err := http.Post(apiURL+"/handshake", "application/json", nil)
	if err != nil || res.StatusCode != 200 {
		log.Printf("Handshake fehlgeschlagen: %v", err)
		return nil, errors.New("Handshake fehlgeschlagen")
	}
	defer res.Body.Close()

	var r handshakeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	cfg = &Config{
		ClientID:  r.ClientID,
		Signature: r.Signature,
	}
	if err := SaveConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
