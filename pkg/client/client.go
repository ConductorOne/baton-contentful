package client

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const defaultBaseURL = "https://api.contentful.com"
const defaultLimit = 100
const contentTypeManagementAPI = "application/vnd.contentful.management.v1+json"

type Client struct {
	*uhttp.BaseHttpClient
	orgID   string
	token   string
	baseURL string
}

func New(ctx context.Context, orgID, token, baseURL string) (*Client, error) {
	client, err := uhttp.NewBearerAuth(token).GetClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	effectiveBaseURL := defaultBaseURL
	if baseURL != "" {
		effectiveBaseURL = baseURL
	}

	return &Client{
		BaseHttpClient: uhttp.NewBaseHttpClient(client),
		orgID:          orgID,
		token:          token,
		baseURL:        effectiveBaseURL,
	}, nil
}
