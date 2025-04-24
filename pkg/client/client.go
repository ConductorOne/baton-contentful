package client

import (
	"net/http"
)

const BaseURL = "https://api.contentful.com"
const defaultLimit = 100

type Client struct {
	http.Client
	orgID string
	token string
}

func NewClient(orgID, token string) *Client {
	return &Client{
		Client: http.Client{},
		orgID:  orgID,
		token:  token,
	}
}
