package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ListSpaces(ctx context.Context, offset int) (*GetSpacesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/spaces", BaseURL), nil)
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

	var res GetSpacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/space-roles
// https://www.contentful.com/help/roles/space-roles-and-permissions/
func (c *Client) ListSpaceRoles(ctx context.Context, spaceID string) (*GetSpaceRolesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/roles", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list users: %s", resp.Status)
	}

	var res GetSpaceRolesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) ListSpaceMemberships(ctx context.Context, spaceID string, offset int) (*GetSpaceMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/spaces/%s/space_memberships", BaseURL, spaceID), nil)
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

	var res GetSpaceMembershipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CreateSpaceMembership(ctx context.Context, spaceID, email string, roleID string, admin bool) (*SpaceMembership, error) {
	body := map[string]interface{}{
		"admin": admin,
		"email": email,
		"roles": []LinkSys{
			{
				Type:     "Link",
				LinkType: "Role",
				ID:       roleID,
			},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/spaces/%s/space_memberships", BaseURL, spaceID), bytes.NewReader(bodyBytes))
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
		return nil, fmt.Errorf("failed to create space membership: %s", resp.Status)
	}

	var res SpaceMembership
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteSpaceMembership(ctx context.Context, spaceID, spaceMembershipID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/spaces/%s/space_memberships/%s", BaseURL, spaceID, spaceMembershipID), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete space membership: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetSpaceMembershipByUser(ctx context.Context, spaceID, userID string) (*GetSpaceMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/space_memberships", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"sys.space.sys.id[eq]": spaceID,
		"sys.user.sys.id[eq]":  userID,
	})

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get space membership by user: %s", resp.Status)
	}

	var res GetSpaceMembershipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &res, nil
}
