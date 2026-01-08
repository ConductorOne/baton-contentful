package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-contentful/pkg/client"
	"github.com/conductorone/baton-sdk/pkg/session"
	"github.com/conductorone/baton-sdk/pkg/types/sessions"
)

type spaceRoleCache struct {
	client *client.Client
}

func newSpaceRoleCache(contentfulClient *client.Client) *spaceRoleCache {
	return &spaceRoleCache{
		client: contentfulClient,
	}
}

type role struct {
	Id   string
	Name string
}

func (c *spaceRoleCache) getSpaceRoles(ctx context.Context, spaceID string) (map[string][]role, error) {
	var offset int
	var spaceRoleCache = make(map[string][]role)
	for {
		res, err := c.client.ListSpaceRoles(ctx, spaceID, offset)
		if err != nil {
			return nil, fmt.Errorf("baton-contentful: failed to list space roles: %w", err)
		}

		if len(res.Items) == 0 {
			break
		}

		for _, r := range res.Items {
			spaceRoleCache[spaceID] = append(spaceRoleCache[spaceID], role{
				Id:   r.Sys.ID,
				Name: r.Name,
			})
		}

		offset += len(res.Items)
	}
	return spaceRoleCache, nil
}

func (c *spaceRoleCache) fillCache(ctx context.Context, spaceID string, sessionStorage sessions.SessionStore) error {
	spaceRoleCache, err := c.getSpaceRoles(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("failed to get space roles: %w", err)
	}

	err = session.SetManyJSON(ctx, sessionStorage, spaceRoleCache, sessions.WithPrefix("contentful_space_role_cache"))
	if err != nil {
		return fmt.Errorf("failed to set space role cache: %w", err)
	}

	return nil
}

func (c *spaceRoleCache) GetCacheRoleNames(ctx context.Context, sessionStorage sessions.SessionStore, spaceID string) (map[string]string, error) {
	spaceRoles, found, err := session.GetJSON[[]role](ctx, sessionStorage, spaceID, sessions.WithPrefix("contentful_space_role_cache"))
	if err != nil {
		return nil, fmt.Errorf("failed to get space role cache: %w", err)
	}
	if !found {
		if err := c.fillCache(ctx, spaceID, sessionStorage); err != nil {
			return nil, fmt.Errorf("failed to fill cache: %w", err)
		}
		// Re-fetch the data from cache after filling it
		spaceRoles, found, err = session.GetJSON[[]role](ctx, sessionStorage, spaceID, sessions.WithPrefix("contentful_space_role_cache"))
		if err != nil {
			return nil, fmt.Errorf("failed to get space role cache after fill: %w", err)
		}
		if !found {
			return nil, fmt.Errorf("space role cache still not found after fill for spaceID: %s", spaceID)
		}
	}

	roleNames := make(map[string]string)
	for _, role := range spaceRoles {
		roleNames[role.Id] = role.Name
	}

	return roleNames, nil
}

// This method uses an in-memory cache instead of sessionStorage.
func (c *spaceRoleCache) GetRoleID(ctx context.Context, spaceID string, roleName string) (string, error) {
	spaceRoles, err := c.getSpaceRoles(ctx, spaceID)
	if err != nil {
		return "", fmt.Errorf("failed to get space roles: %w", err)
	}

	for _, role := range spaceRoles[spaceID] {
		if role.Name == roleName {
			return role.Id, nil
		}
	}
	return "", fmt.Errorf("role %s not found in cache, spaceID: %s", roleName, spaceID)
}
