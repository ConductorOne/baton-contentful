package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ListTeams(ctx context.Context, offset int) (*GetTeamsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/teams", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	url := req.URL.Query()
	url.Set("limit", fmt.Sprintf("%d", defaultLimit))
	url.Set("skip", fmt.Sprintf("%d", offset))
	req.URL.RawQuery = url.Encode()

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list users: %s", resp.Status)
	}

	var res GetTeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) ListTeamMemberships(ctx context.Context, offset int) (*GetTeamMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/team_memberships", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	url := req.URL.Query()
	url.Set("limit", fmt.Sprintf("%d", defaultLimit))
	url.Set("skip", fmt.Sprintf("%d", offset))
	req.URL.RawQuery = url.Encode()

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list users: %s", resp.Status)
	}

	var res GetTeamMembershipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
