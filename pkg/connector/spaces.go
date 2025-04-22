package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-contentful/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type spaceBuilder struct {
	client *client.Client
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

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
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
	users := res.Items
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(users) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(users))

	rv := make([]*v2.Resource, len(users))
	for i, user := range users {
		rv[i] = spaceResource(user)
	}

	return rv, nextOffset, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *spaceBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *spaceBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newSpaceBuilder(client *client.Client) *spaceBuilder {
	return &spaceBuilder{
		client: client,
	}
}
