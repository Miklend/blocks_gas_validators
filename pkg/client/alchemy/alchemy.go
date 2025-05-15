package alchemyClient

import (
	"blocks_gas_validators/internal/configs"
	"blocks_gas_validators/pkg/chains"
	"blocks_gas_validators/pkg/logging"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type Client struct {
	NetworkName string
	APIKey      string
	BaseURL     string
	Client      *ethclient.Client
}

func NewAlchemyClient(cfg configs.AlchemyConfig, logger *logging.Logger) (*Client, error) {
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("error loading env variables: %s", err.Error())
	}
	apiKey := os.Getenv(cfg.NameApiKey)

	var fullURL string
	if cfg.Mode == "last" {
		fullURL = fmt.Sprintf("wss%s%s", chains.AlchemyChains[cfg.NetworkName].URL, apiKey)
	} else {
		fullURL = fmt.Sprintf("https%s%s", chains.AlchemyChains[cfg.NetworkName].URL, apiKey)
	}

	client, err := ethclient.Dial(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed connect to %s: %w", cfg.NetworkName, err)
	}

	return &Client{
		NetworkName: cfg.NetworkName,
		APIKey:      apiKey,
		BaseURL:     chains.AlchemyChains[cfg.NetworkName].URL,
		Client:      client,
	}, nil
}

func (a *Client) Close() {
	a.Client.Close()
}
