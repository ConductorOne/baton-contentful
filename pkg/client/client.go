package client

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const BaseURL = "https://api.contentful.com"
const defaultLimit = 100

type Client struct {
	*uhttp.BaseHttpClient
	orgID string
	token string
}

func New(ctx context.Context, orgID, token string) (*Client, error) {
	client, err := uhttp.NewBearerAuth(token).GetClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}
	return &Client{
		BaseHttpClient: uhttp.NewBaseHttpClient(client),
		orgID:          orgID,
		token:          token,
	}, nil
}
