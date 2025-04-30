package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListUsers(ctx context.Context, offset int) (*GetUsersResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/users", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetUsersResponse
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list users: %s", resp.Status)
	}

	return &res, nil
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (*GetUsersResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/users", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"query": userID,
	})

	var res GetUsersResponse
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user by id: %s", resp.Status)
	}

	return &res, nil
}

func (c *Client) CreateInvitation(ctx context.Context, body *CreateInvitationBody) (*Invitation, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/%s/invitations", BaseURL, c.orgID), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.contentful.management.v1+json")

	var res Invitation
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create space membership: %s", resp.Status)
	}

	return &res, nil
}

func (c *Client) GetLastActiveAt(ctx context.Context, userID string) *time.Time {
	res, err := c.GetOrganizationMembershipByUser(ctx, userID)
	if err != nil {
		return nil
	}

	if len(res.Items) == 0 {
		return nil
	}

	return res.Items[0].Sys.LastActiveAt
}
