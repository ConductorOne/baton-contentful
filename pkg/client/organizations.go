package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListOrganizations(ctx context.Context, offset int) (*GetOrganizationsResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "organizations")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetOrganizationsResponse
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

// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/organization-memberships
func (c *Client) ListOrganizationMemberships(ctx context.Context, offset int) (*GetOrganizationMembershipsResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "organizations", c.orgID, "organization_memberships")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"limit": fmt.Sprintf("%d", defaultLimit),
		"skip":  fmt.Sprintf("%d", offset),
	})

	var res GetOrganizationMembershipsResponse
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

func (c *Client) GetOrganizationMembershipByUser(ctx context.Context, userID string) (*GetOrganizationMembershipsResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "organizations", c.orgID, "organization_memberships")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	SetQueryParams(req.URL, map[string]string{
		"sys.user.sys.id[eq]": userID,
	})

	var res GetOrganizationMembershipsResponse
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

func (c *Client) DeleteOrganizationMembership(ctx context.Context, orgMembershipID string) error {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "organizations", c.orgID, "organization_memberships", orgMembershipID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL.String(), nil)
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
