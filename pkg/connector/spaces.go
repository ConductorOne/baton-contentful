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
	// spaceId: role
	spaceRoleCache map[string][]role
	mu             *sync.Mutex
}

type role struct {
	Id   string
	Name string
}

func (o *spaceBuilder) fillCache(ctx context.Context, spaceID string) error {
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
			o.cacheSetRole(spaceID, role.Sys.ID, role.Name)
		}

		offset += len(res.Items)
	}
	return nil
}

func (o *spaceBuilder) cacheGetRoleName(ctx context.Context, spaceID, roleID string) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// if no roles are cached for the space, we need to fill the cache
	if len(o.spaceRoleCache[spaceID]) == 0 {
		if err := o.fillCache(ctx, spaceID); err != nil {
			return "", fmt.Errorf("failed to fill cache: %w", err)
		}
	}

	for _, role := range o.spaceRoleCache[spaceID] {
		if role.Id == roleID {
			return role.Name, nil
		}
	}

	return "", fmt.Errorf("roleID %s not found in cache, spaceID: %s", roleID, spaceID)
}

func (o *spaceBuilder) cacheGetRoleID(ctx context.Context, spaceID, roleName string) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// if no roles are cached for the space, we need to fill the cache
	if len(o.spaceRoleCache[spaceID]) == 0 {
		if err := o.fillCache(ctx, spaceID); err != nil {
			return "", fmt.Errorf("failed to fill cache: %w", err)
		}
	}

	for _, role := range o.spaceRoleCache[spaceID] {
		if role.Name == roleName {
			return role.Id, nil
		}
	}
	return "", fmt.Errorf("role %s not found in cache, spaceID: %s", roleName, spaceID)
}

func (o *spaceBuilder) cacheSetRole(spaceId, roleID, roleName string) {
	o.spaceRoleCache[spaceId] = append(o.spaceRoleCache[spaceId], role{
		Id:   roleID,
		Name: roleName,
	})
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

	rv := []*v2.Entitlement{}
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

	res, err := o.client.ListSpaceRoles(ctx, resource.Id.Resource, offset)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list space roles: %w", err)
	}

	if len(res.Items) == 0 {
		return rv, "", nil, nil
	}

	for _, role := range res.Items {
		rv = append(rv, entitlement.NewAssignmentEntitlement(
			resource,
			role.Name,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Role %s for %s space", role.Name, resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Role %s for %s space ", role.Name, resource.DisplayName)),
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

	// admin role is special, we don't need to look it up
	if roleName != spaceAdmin {
		roleID, err = o.cacheGetRoleID(ctx, spaceID, roleName)
		if err != nil {
			return nil, fmt.Errorf("baton-contentful: failed to get role ID for role %s: %w", roleName, err)
		}
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
		client:         client,
		mu:             &sync.Mutex{},
		spaceRoleCache: make(map[string][]role),
	}
}
