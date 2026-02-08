package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListUsers(ctx context.Context, offset int) (*GetUsersResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "users")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
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

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (*GetUsersResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "users")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
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

	defer resp.Body.Close()

	return &res, nil
}

func (c *Client) CreateInvitation(ctx context.Context, body *CreateInvitationBody) (*Invitation, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	reqURL.Path = path.Join(reqURL.Path, "organizations", c.orgID, "invitations")

	req, err := c.NewRequest(ctx, http.MethodPost, reqURL,
		uhttp.WithJSONBody(body),
		uhttp.WithHeader("Content-Type", contentTypeManagementAPI),
	)
	if err != nil {
		return nil, err
	}

	var res Invitation
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
