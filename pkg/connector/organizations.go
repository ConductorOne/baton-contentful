package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-contentful/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	orgOwner     = "owner"
	orgAdmin     = "admin"
	orgDeveloper = "developer"
	orgMember    = "member"
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
		return nil, "", nil, fmt.Errorf("baton-contentful: failed to list users: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))

	rv := []*v2.Resource{}
	for _, org := range res.Items {
		rv = append(rv, orgResource(org))
	}

	return rv, nextOffset, nil, nil
}

func (o *orgBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	// owner, admin, developer, member
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			orgOwner,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Owner of the %s organization", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Owner of the %s organization", resource.DisplayName)),
		),
		entitlement.NewAssignmentEntitlement(
			resource,
			orgAdmin,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Admin of the %s organization", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Admin of the %s organization", resource.DisplayName)),
		),
		entitlement.NewAssignmentEntitlement(
			resource,
			orgDeveloper,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Developer of the %s organization", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Developer of the %s organization", resource.DisplayName)),
		),
		entitlement.NewAssignmentEntitlement(
			resource,
			orgMember,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Member of the %s organization", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of the %s organization", resource.DisplayName)),
		),
	}, "", nil, nil
}

func (o *orgBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListOrganizationMemberships(ctx, offset)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-contentful: failed to list org memberships: %w", err)
	}

	if len(res.Items) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(res.Items))

	rv := []*v2.Grant{}
	for _, orgMembership := range res.Items {
		principalID, err := resourceSdk.NewResourceID(userResourceType, orgMembership.Sys.User.Sys.ID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-contentful: failed to create resource ID for user %v: %w", orgMembership.Sys.User.Sys.ID, err)
		}
		rv = append(rv, grant.NewGrant(
			resource,
			orgMembership.Role,
			principalID,
		))
	}
	return rv, nextOffset, nil, nil
}

// can't provision organization membership, it requires creating an account
// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/organization-memberships
func (o *orgBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	return nil, nil
}

func (o *orgBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	principal := grant.Principal

	resOrgMembership, err := o.client.GetOrganizationMembershipByUser(ctx, principal.Id.Resource)
	if err != nil {
		return nil, err
	}

	if len(resOrgMembership.Items) == 0 {
		return annotations.New(&v2.GrantAlreadyRevoked{}), fmt.Errorf("baton-contentful: organization membership not found for user %s", principal.Id.Resource)
	}

	orgMembershipID := resOrgMembership.Items[0].Sys.ID
	err = o.client.DeleteOrganizationMembership(ctx, orgMembershipID)
	if err != nil {
		return nil, fmt.Errorf("baton-contentful: failed to delete organization membership %s: %w", orgMembershipID, err)
	}
	return nil, nil
}

func newOrgBuilder(client *client.Client) *orgBuilder {
	return &orgBuilder{
		client: client,
	}
}
