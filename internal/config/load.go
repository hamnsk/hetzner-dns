package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type Zones struct {
	Zones []struct {
		Name    string `yaml:"name"`
		TTL     int    `yaml:"ttl"`
		Records []struct {
			Name  string `yaml:"name"`
			Type  string `yaml:"type"`
			Value string `yaml:"value,omitempty"`
			TTL   int    `yaml:"ttl"`
		} `yaml:"records"`
	} `yaml:"zones"`
}

type AppConfig struct {
	HetznerToken         string
	ZonesConfigFile      string
	UpdateInterval       string
	UseVault             bool
	VaultRoleID          string
	VaultAddr            string
	VaultAuthName        string
	VaultSecretStoreName string
	VaultSecretStorePath string
}

func GetAppConfig() AppConfig {
	flag, _ := strconv.ParseBool(os.Getenv("HETZNER_DNS_USE_VAULT"))
	return AppConfig{
		HetznerToken:         os.Getenv("HETZNER_DNS_TOKEN"),
		ZonesConfigFile:      os.Getenv("HETZNER_DNS_ZONES_CONFIG_FILE"),
		UpdateInterval:       os.Getenv("HETZNER_DNS_ZONES_UPDATE_INTERVAL"),
		UseVault:             flag,
		VaultAddr:            os.Getenv("HETZNER_DNS_VAULT_ADDR"),
		VaultAuthName:        os.Getenv("HETZNER_DNS_VAULT_AUTH_NAME"),
		VaultRoleID:          os.Getenv("HETZNER_DNS_VAULT_ROLE_ID"),
		VaultSecretStoreName: os.Getenv("HETZNER_DNS_VAULT_SECRET_STORE_NAME"),
		VaultSecretStorePath: os.Getenv("HETZNER_DNS_VAULT_SECRET_STORE_PATH"),
	}
}

func (a *AppConfig) GetZones() (Zones, error) {

	confFile, err := os.ReadFile(a.ZonesConfigFile)
	if err != nil {
		return Zones{}, err
	}

	var zones Zones
	err = yaml.Unmarshal(confFile, &zones)
	if err != nil {
		return Zones{}, err
	}

	return zones, nil
}
