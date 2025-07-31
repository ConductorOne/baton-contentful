package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
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

	var res GetSpacesResponse
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

// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/space-roles
// https://www.contentful.com/help/roles/space-roles-and-permissions/
func (c *Client) ListSpaceRoles(ctx context.Context, spaceID string, offset int) (*GetSpaceRolesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/%s/roles", BaseURL, c.orgID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetSpaceRolesResponse
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

func (c *Client) ListSpaceMemberships(ctx context.Context, spaceID string, offset int) (*GetSpaceMembershipsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/spaces/%s/space_memberships", BaseURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetSpaceMembershipsResponse
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

func (c *Client) CreateSpaceMembership(ctx context.Context, spaceID, email string, roleID string, isAdmin bool) (*SpaceMembership, error) {
	body := map[string]interface{}{
		"admin": isAdmin,
		"email": email,
	}

	if roleID != "" {
		body["roles"] = []LinkSys{
			{
				Type:     "Link",
				LinkType: "Role",
				ID:       roleID,
			},
		}
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

	var res SpaceMembership
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

func (c *Client) DeleteSpaceMembership(ctx context.Context, spaceID, spaceMembershipID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/spaces/%s/space_memberships/%s", BaseURL, spaceID, spaceMembershipID), nil)
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

	var res GetSpaceMembershipsResponse
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
