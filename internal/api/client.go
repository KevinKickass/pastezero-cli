package api

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	ClientID  string `json:"client_id"`
	Signature string `json:"signature"`
}

func configPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "pastezero", "config.json")
}

func LoadConfig() (*Config, error) {
	f, err := os.Open(configPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	if c.ClientID == "" || c.Signature == "" {
		return nil, errors.New("ungültige Konfiguration")
	}
	return &c, nil
}

func SaveConfig(c *Config) error {
	dir := filepath.Dir(configPath())
	_ = os.MkdirAll(dir, 0700)
	f, err := os.Create(configPath())
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c)
}
