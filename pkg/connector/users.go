package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-contentful/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client *client.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func (o *userBuilder) userResource(ctx context.Context, user client.User) *v2.Resource {
	profile := map[string]interface{}{
		"firstName":  user.FirstName,
		"lastName":   user.LastName,
		"email":      user.Email,
		"2faEnabled": user.TwoFAEnabled,
	}

	traits := []resourceSdk.UserTraitOption{
		resourceSdk.WithEmail(user.Email, true),
		resourceSdk.WithUserProfile(profile),
		resourceSdk.WithCreatedAt(user.Sys.CreatedAt),
	}

	lastActive := o.client.GetLastActiveAt(ctx, user.Sys.ID)
	if lastActive != nil {
		traits = append(traits, resourceSdk.WithLastLogin(*lastActive))

	}

	userResource, err := resourceSdk.NewUserResource(
		fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		userResourceType,
		user.Sys.ID,
		traits,
	)
	if err != nil {
		return nil
	}

	return userResource
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var offset int
	var err error
	if pToken.Token != "" {
		offset, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, err
		}
	}

	res, err := o.client.ListUsers(ctx, offset)
	users := res.Items
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(users) == 0 {
		return nil, "", nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(users))

	rv := []*v2.Resource{}
	for _, user := range users {
		rv = append(rv, o.userResource(ctx, user))
	}

	return rv, nextOffset, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *userBuilder) CreateAccountCapabilityDetails(ctx context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

func (o *userBuilder) CreateAccount(ctx context.Context, accountInfo *v2.AccountInfo, credentialOptions *v2.CredentialOptions) (
	connectorbuilder.CreateAccountResponse,
	[]*v2.PlaintextData,
	annotations.Annotations,
	error,
) {
	body, err := getCreateInvitationBody(accountInfo)
	if err != nil {
		return nil, nil, nil, err
	}

	invitation, err := o.client.CreateInvitation(ctx, body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("baton-contentful: cannot create invitation: %w", err)
	}

	return &v2.CreateAccountResponse_ActionRequiredResult{
		Message: invitation.Sys.InvitationURL,
	}, nil, nil, nil
}

func getCreateInvitationBody(accountInfo *v2.AccountInfo) (*client.CreateInvitationBody, error) {
	pMap := accountInfo.Profile.AsMap()
	firstName := pMap["firstName"].(string)
	lastName := pMap["lastName"].(string)
	role := pMap["role"].(string)

	return &client.CreateInvitationBody{
		Email:     accountInfo.Login,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
