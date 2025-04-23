package connector

import (
	"context"
	"fmt"
	"strconv"
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
	// id: name
	roleCache map[string]string
	mu        *sync.Mutex
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
		return nil, "", nil, fmt.Errorf("failed to list users: %w", err)
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
		return nil, "", nil, fmt.Errorf("failed to list org memberships: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))

	rv := []*v2.Grant{}
	for _, spaceMembership := range res.Items {
		principalID, err := resourceSdk.NewResourceID(userResourceType, spaceMembership.Sys.User.Sys.ID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create resource ID for user %v: %w", spaceMembership.Sys.User.Sys.ID, err)
		}

		if spaceMembership.Admin {
			principalID, err = resourceSdk.NewResourceID(userResourceType, spaceMembership.Sys.User.Sys.ID)
			if err != nil {
				return nil, "", nil, fmt.Errorf("failed to create resource ID for user %v: %w", spaceMembership.Sys.User.Sys.ID, err)
			}
			rv = append(rv, grant.NewGrant(
				resource,
				spaceAdmin,
				principalID,
			))
			continue
		}

		for _, role := range spaceMembership.Roles {
			roleName := o.cacheGetRole(role.Sys.ID)
			if roleName == "" {
				logger.Info("cache miss for role", zap.String("roleID", role.Sys.ID))
				return nil, "", nil, fmt.Errorf("failed to get role name for role %s", role.Sys.ID)
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

func newSpaceBuilder(client *client.Client) *spaceBuilder {
	return &spaceBuilder{
		client:    client,
		mu:        &sync.Mutex{},
		roleCache: make(map[string]string),
	}
}
