package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListSpaces(ctx context.Context, offset int) (*GetSpacesResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "spaces")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
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

// https://www.contentful.com/developers/docs/references/content-management-api/#/reference/roles/roles-collection/get-all-roles/console/curl
// https://www.contentful.com/help/roles/space-roles-and-permissions/
func (c *Client) ListSpaceRoles(ctx context.Context, spaceID string, offset int) (*GetSpaceRolesResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "spaces", spaceID, "roles")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
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

func (c *Client) ListSpaceMembers(ctx context.Context, spaceID string, offset int) (*GetSpaceMembershipsResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "spaces", spaceID, "space_members")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
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

	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "spaces", spaceID, "space_memberships")
	reqURL := baseURL

	req, err := c.NewRequest(ctx, http.MethodPost, reqURL,
		uhttp.WithJSONBody(body),
		uhttp.WithHeader("Content-Type", "application/vnd.contentful.management.v1+json"),
	)
	if err != nil {
		return nil, err
	}

	var res SpaceMembership
	var errorResp ErrorResponse
	resp, err := c.Do(req,
		uhttp.WithJSONResponse(&res),
		uhttp.WithErrorResponse(&errorResp),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) UpdateSpaceMembership(ctx context.Context, spaceID, spaceMembershipID, email string, roles []LinkSys, isAdmin bool, version int) error {
	body := map[string]interface{}{
		"admin": isAdmin,
		"email": email,
	}

	if roles != nil {
		body["roles"] = roles
	}

	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "spaces", spaceID, "space_memberships", spaceMembershipID)
	reqURL := baseURL

	req, err := c.NewRequest(ctx, http.MethodPut, reqURL,
		uhttp.WithJSONBody(body),
		uhttp.WithHeader("Content-Type", "application/vnd.contentful.management.v1+json"),
		uhttp.WithHeader("X-Contentful-Version", fmt.Sprintf("%d", version)),
	)
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

func (c *Client) DeleteSpaceMembership(ctx context.Context, spaceID, spaceMembershipID string) error {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "spaces", spaceID, "space_memberships", spaceMembershipID)

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

func (c *Client) GetSpaceMembershipByUser(ctx context.Context, spaceID, userID string) (*GetSpaceMembershipsResponse, error) {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "organizations", c.orgID, "space_memberships")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
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
