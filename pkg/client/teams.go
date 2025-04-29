package client

import (
	"bytes"
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

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	resp, err := c.Do(req)
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

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	resp, err := c.Do(req)
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

func (c *Client) CreateTeamMembership(ctx context.Context, teamID string, orgMembershipID string) (*TeamMembership, error) {
	body := map[string]interface{}{
		"organizationMembershipId": orgMembershipID,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/%s/teams/%s/team_memberships", BaseURL, c.orgID, teamID), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/vnd.contentful.management.v1+json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create team membership: %s", resp.Status)
	}

	var res TeamMembership
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTeamMembershipByUser(ctx context.Context, orgMembershipID string) (*GetTeamMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/team_memberships", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"sys.organizationMembership.sys.id": orgMembershipID,
	})

	resp, err := c.Do(req)
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

func (c *Client) DeleteTeamMembership(ctx context.Context, teamID, teamMembershipID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/%s/teams/%s/team_memberships/%s", BaseURL, c.orgID, teamID, teamMembershipID), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete team membership: %s", resp.Status)
	}

	return nil
}
