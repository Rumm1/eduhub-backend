package redis

import "context"

type Config struct {
	Address  string
	Password string
	DB       int
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}

func (c *Client) Ping(_ context.Context) error {
	return nil
}

func (c *Client) Close() error {
	return nil
}
