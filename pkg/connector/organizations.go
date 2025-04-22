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

type orgBuilder struct {
	client *client.Client
}

func (o *orgBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return orgResourceType
}

func orgResource(org client.Organization) *v2.Resource {
	orgResource, err := resourceSdk.NewGroupResource(
		org.Name,
		orgResourceType,
		org.Sys.ID,
		[]resourceSdk.GroupTraitOption{},
	)
	if err != nil {
		return nil
	}

	return orgResource
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *orgBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListOrganizations(ctx, offset)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))

	rv := make([]*v2.Resource, len(res.Items))
	for i, user := range res.Items {
		rv[i] = orgResource(user)
	}

	return rv, nextOffset, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *orgBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *orgBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newOrgBuilder(client *client.Client) *orgBuilder {
	return &orgBuilder{
		client: client,
	}
}
