package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListTeams(ctx context.Context, offset int) (*GetTeamsResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "teams")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetTeamsResponse
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) ListTeamMembershipsByTeam(ctx context.Context, teamID string, offset int) (*GetTeamMembershipsResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "teams", teamID, "team_memberships")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetTeamMembershipsResponse
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) CreateTeamMembership(ctx context.Context, teamID string, orgMembershipID string) (*TeamMembership, error) {
	body := map[string]interface{}{
		"organizationMembershipId": orgMembershipID,
	}

	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "teams", teamID, "team_memberships")

	req, err := c.NewRequest(ctx, http.MethodPost, reqURL,
		uhttp.WithJSONBody(body),
		uhttp.WithHeader("Content-Type", contentTypeManagementAPI),
	)
	if err != nil {
		return nil, err
	}

	var res TeamMembership
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) GetTeamMembershipByUser(ctx context.Context, orgMembershipID string) (*GetTeamMembershipsResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "team_memberships")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"sys.organizationMembership.sys.id": orgMembershipID,
	})

	var res GetTeamMembershipsResponse
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) DeleteTeamMembership(ctx context.Context, teamID, teamMembershipID string) error {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "teams", teamID, "team_memberships", teamMembershipID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req,
		uhttp.WithErrorResponse(&ErrorResponse{}),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
