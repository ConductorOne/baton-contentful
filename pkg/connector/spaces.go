package connector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/conductorone/baton-contentful/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const spaceAdmin = "admin"

type spaceBuilder struct {
	client *client.Client
	// roleID: name
	roleCache map[string]string
	mu        *sync.Mutex
	once      *sync.Once
}

func (o spaceBuilder) fillCache(ctx context.Context, spaceID string) error {
	var offset int
	for {
		res, err := o.client.ListSpaceRoles(ctx, spaceID, offset)
		if err != nil {
			return fmt.Errorf("baton-contentful: failed to list space roles: %w", err)
		}

		if len(res.Items) == 0 {
			break
		}

		for _, role := range res.Items {
			o.cacheSetRole(role.Sys.ID, role.Name)
		}

		offset += len(res.Items)
	}
	return nil
}

func (o spaceBuilder) cacheGetRoleName(ctx context.Context, spaceID, roleID string) (string, error) {
	var fillErr error
	o.once.Do(func() {
		fillErr = o.fillCache(ctx, spaceID)
	})
	if fillErr != nil {
		return "", fmt.Errorf("failed to fill cache: %w", fillErr)
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if role, ok := o.roleCache[roleID]; ok {
		return role, nil
	}

	return "", fmt.Errorf("roleID %s not found in cache, spaceID: %s", roleID, spaceID)
}

func (o spaceBuilder) cacheGetRoleID(ctx context.Context, spaceID, roleName string) (string, error) {
	var fillErr error
	o.once.Do(func() {
		fillErr = o.fillCache(ctx, spaceID)
	})
	if fillErr != nil {
		return "", fmt.Errorf("failed to fill cache: %w", fillErr)
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	for roleID, name := range o.roleCache {
		if name == roleName {
			return roleID, nil
		}
	}
	return "", fmt.Errorf("role %s not found in cache, spaceID: %s", roleName, spaceID)
}

func (o *spaceBuilder) cacheSetRole(roleID string, roleName string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.roleCache[roleID] = roleName
}

func (o *spaceBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return spaceResourceType
}

func spaceResource(space client.Space) *v2.Resource {
	spaceResource, err := resourceSdk.NewGroupResource(
		space.Name,
		spaceResourceType,
		space.Sys.ID,
		[]resourceSdk.GroupTraitOption{},
	)
	if err != nil {
		return nil
	}

	return spaceResource
}

func (o *spaceBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListSpaces(ctx, offset)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-contentful: failed to list users: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))

	rv := []*v2.Resource{}
	for _, space := range res.Items {
		rv = append(rv, spaceResource(space))
	}

	return rv, nextOffset, nil, nil
}

func (o *spaceBuilder) Entitlements(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}
	res, err := o.client.ListSpaceRoles(ctx, resource.Id.Resource, offset)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list space roles: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}

	rv := make([]*v2.Entitlement, 0, len(res.Items))
	for _, role := range res.Items {
		o.cacheSetRole(role.Sys.ID, role.Name)
		rv = append(rv, entitlement.NewAssignmentEntitlement(
			resource,
			role.Name,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Role %s for %s space", role.Name, resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Role %s for %s space ", role.Name, resource.DisplayName)),
		))
	}

	// so it's only included once
	if offset == 0 {
		rv = append(rv, entitlement.NewAssignmentEntitlement(
			resource,
			spaceAdmin,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Admin for %s space", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Admin for %s space", resource.DisplayName)),
		))
	}

	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))
	return rv, nextOffset, nil, nil
}

func (o *spaceBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListSpaceMembers(ctx, resource.Id.Resource, offset)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-contentful: failed to list org memberships: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))

	rv := []*v2.Grant{}
	for _, spaceMembership := range res.Items {
		principalID, err := resourceSdk.NewResourceID(userResourceType, spaceMembership.Sys.User.Sys.ID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-contentful: failed to create resource ID for user %v: %w", spaceMembership.Sys.User.Sys.ID, err)
		}

		if spaceMembership.Admin {
			principalID, err = resourceSdk.NewResourceID(userResourceType, spaceMembership.Sys.User.Sys.ID)
			if err != nil {
				return nil, "", nil, fmt.Errorf("baton-contentful: failed to create resource ID for user %v: %w", spaceMembership.Sys.User.Sys.ID, err)
			}
			rv = append(rv, grant.NewGrant(
				resource,
				spaceAdmin,
				principalID,
			))
			continue
		}

		for _, role := range spaceMembership.Roles {
			roleName, err := o.cacheGetRoleName(ctx, resource.Id.Resource, role.Sys.ID)
			if err != nil {
				return nil, "", nil, fmt.Errorf("baton-contentful: failed to get role name for role ID %s: %w", role.Sys.ID, err)
			}
			rv = append(rv, grant.NewGrant(
				resource,
				roleName,
				principalID,
			))
		}
	}
	return rv, nextOffset, nil, nil
}

func (o *spaceBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	spaceID := entitlement.Resource.Id.Resource
	roleName := strings.Split(entitlement.Id, ":")[2]

	resUser, err := o.client.GetUserByID(ctx, principal.Id.Resource)
	if err != nil {
		return nil, err
	}
	if len(resUser.Items) == 0 {
		return nil, fmt.Errorf("baton-contentful: no user found for ID %s", principal.Id.Resource)
	}

	roleID := ""
	isAdmin := roleName == spaceAdmin

	roleID, err = o.cacheGetRoleID(ctx, spaceID, roleName)
	if err != nil {
		return nil, fmt.Errorf("baton-contentful: failed to get role ID for role %s: %w", roleName, err)
	}

	email := resUser.Items[0].Email

	// if the user is not an admin and no role ID is provided, we cannot create the membership
	// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/space-memberships
	if !isAdmin && roleID == "" {
		return nil, fmt.Errorf("baton-contentful: role ID must be provided for non-admin space memberships")
	}

	_, err = o.client.CreateSpaceMembership(ctx, spaceID, email, roleID, isAdmin)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (o *spaceBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	principal := grant.Principal
	entitlement := grant.Entitlement
	spaceID := entitlement.Resource.Id.Resource

	resSpaceMembership, err := o.client.GetSpaceMembershipByUser(ctx, spaceID, principal.Id.Resource)
	if err != nil {
		return nil, err
	}

	if len(resSpaceMembership.Items) == 0 {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	spaceMembershipID := resSpaceMembership.Items[0].Sys.ID
	err = o.client.DeleteSpaceMembership(ctx, spaceID, spaceMembershipID)
	if err != nil {
		return nil, fmt.Errorf("baton-contentful: failed to delete team membership: %w", err)
	}
	return nil, nil
}

func newSpaceBuilder(client *client.Client) *spaceBuilder {
	return &spaceBuilder{
		client:    client,
		mu:        &sync.Mutex{},
		roleCache: make(map[string]string),
		once:      &sync.Once{},
	}
}
