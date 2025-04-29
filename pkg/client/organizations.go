package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ListOrganizations(ctx context.Context, offset int) (*GetOrganizationsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations", BaseURL), nil)
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

	var res GetOrganizationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/organization-memberships
func (c *Client) ListOrganizationMemberships(ctx context.Context, offset int) (*GetOrganizationMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/organization_memberships", BaseURL, c.orgID), nil)
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

	var res GetOrganizationMembershipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetOrganizationMembershipByUser(ctx context.Context, userID string) (*GetOrganizationMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/organization_memberships", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"sys.user.sys.id[eq]": userID,
	})

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get organization membership: %s", resp.Status)
	}

	var res GetOrganizationMembershipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteOrganizationMembership(ctx context.Context, orgMembershipID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/%s/organization_memberships/%s", BaseURL, c.orgID, orgMembershipID), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete organization membership: %s", resp.Status)
	}

	return nil
}
