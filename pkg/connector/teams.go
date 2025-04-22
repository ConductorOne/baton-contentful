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

type teamBuilder struct {
	client *client.Client
}

func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

func teamResource(team client.Team) *v2.Resource {
	teamResource, err := resourceSdk.NewGroupResource(
		team.Name,
		teamResourceType,
		team.Sys.ID,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(
				map[string]interface{}{
					"description": team.Description,
				},
			),
		},
	)
	if err != nil {
		return nil
	}

	return teamResource
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *teamBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListTeams(ctx, offset)
	items := res.Items
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(items))

	rv := make([]*v2.Resource, len(items))
	for i, elem := range items {
		rv[i] = teamResource(elem)
	}

	return rv, nextOffset, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newTeamBuilder(client *client.Client) *teamBuilder {
	return &teamBuilder{
		client: client,
	}
}
