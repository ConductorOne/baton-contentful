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
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const spaceAdmin = "admin"

type spaceBuilder struct {
	client *client.Client
	// roleID: name
	roleCache map[string]string
	mu        *sync.Mutex
}

func (o spaceBuilder) fillCache(ctx context.Context, spaceID string) error {
	res, err := o.client.ListSpaceRoles(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("baton-contentful: failed to list space roles: %w", err)
	}
	for _, role := range res.Items {
		o.cacheSetRole(role.Sys.ID, role.Name)
	}
	return nil
}

func (o spaceBuilder) cacheGetRole(roleID string) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if role, ok := o.roleCache[roleID]; ok {
		return role
	}

	return ""
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

func (o *spaceBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	res, err := o.client.ListSpaceRoles(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list space roles: %w", err)
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

	rv = append(rv, entitlement.NewAssignmentEntitlement(
		resource,
		spaceAdmin,
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Admin for %s space", resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("Admin for %s space", resource.DisplayName)),
	))

	return rv, "", nil, nil
}

func (o *spaceBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	logger := ctxzap.Extract(ctx)
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListSpaceMemberships(ctx, resource.Id.Resource, offset)
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
			if len(o.roleCache) == 0 {
				err := o.fillCache(ctx, resource.Id.Resource)
				if err != nil {
					return nil, "", nil, fmt.Errorf("baton-contentful: failed to fill cache: %w", err)
				}
			}

			roleName := o.cacheGetRole(role.Sys.ID)
			if roleName == "" {
				logger.Info("cache miss for role", zap.String("roleID", role.Sys.ID))
				return nil, "", nil, fmt.Errorf("baton-contentful: failed to get role name for role %s", role.Sys.ID)
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
	role := strings.Split(entitlement.Id, ":")[2]

	resUser, err := o.client.GetUserByID(ctx, principal.Id.Resource)
	if err != nil {
		return nil, err
	}
	if len(resUser.Items) == 0 {
		return nil, fmt.Errorf("baton-contentful: no user found for ID %s", principal.Id.Resource)
	}

	roleID := ""
	admin := role == spaceAdmin
	if !admin {
		resSpaceRoles, err := o.client.ListSpaceRoles(ctx, spaceID)
		if err != nil {
			return nil, fmt.Errorf("baton-contentful: failed to list space roles: %w", err)
		}

		for _, item := range resSpaceRoles.Items {
			if item.Name == role {
				roleID = item.Sys.ID
				break
			}
		}
	}

	email := resUser.Items[0].Email
	_, err = o.client.CreateSpaceMembership(ctx, spaceID, email, roleID, admin)
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
	}
}
