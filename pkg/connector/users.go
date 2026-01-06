package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-contentful/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client *client.Client
}

var _ connectorbuilder.AccountManagerV2 = &userBuilder{}

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

	traits := []rs.UserTraitOption{
		rs.WithEmail(user.Email, true),
		rs.WithUserProfile(profile),
		rs.WithCreatedAt(user.Sys.CreatedAt),
	}

	lastActive := o.client.GetLastActiveAt(ctx, user.Sys.ID)
	if lastActive != nil {
		traits = append(traits, rs.WithLastLogin(*lastActive))
	}

	userResource, err := rs.NewUserResource(
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
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, attrs rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var offset int
	var err error
	if attrs.PageToken.Token != "" {
		offset, err = strconv.Atoi(attrs.PageToken.Token)
		if err != nil {
			return nil, nil, err
		}
	}

	res, err := o.client.ListUsers(ctx, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("baton-contentful: failed to list users: %w", err)
	}

	users := res.Items
	if len(users) == 0 {
		return nil, nil, nil
	}
	nextOffset := fmt.Sprintf("%d", offset+len(users))

	rv := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		rv = append(rv, o.userResource(ctx, user))
	}

	return rv, &rs.SyncOpResults{NextPageToken: nextOffset}, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

func (o *userBuilder) CreateAccountCapabilityDetails(ctx context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

func (o *userBuilder) CreateAccount(ctx context.Context, accountInfo *v2.AccountInfo, credentialOptions *v2.LocalCredentialOptions) (
	connectorbuilder.CreateAccountResponse,
	[]*v2.PlaintextData,
	annotations.Annotations,
	error,
) {
	body, err := getCreateInvitationBody(accountInfo)
	if err != nil {
		return nil, nil, nil, err
	}

	_, err = o.client.CreateInvitation(ctx, body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("baton-contentful: cannot create invitation: %w", err)
	}

	var userResource *v2.Resource
	resUser, err := o.client.GetUserByID(ctx, body.Email)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("baton-contentful: failed to get user after invitation: %w", err)
	}
	if len(resUser.Items) > 0 {
		userResource = o.userResource(ctx, resUser.Items[0])
	}

	return &v2.CreateAccountResponse_SuccessResult{
		Resource:              userResource,
		IsCreateAccountResult: true,
	}, nil, nil, nil
}

func getCreateInvitationBody(accountInfo *v2.AccountInfo) (*client.CreateInvitationBody, error) {
	pMap := accountInfo.Profile.AsMap()
	firstName := ""
	lastName := ""
	role := ""
	email := ""

	if pMap["firstName"] != nil {
		firstName = pMap["firstName"].(string)
	}
	if pMap["lastName"] != nil {
		lastName = pMap["lastName"].(string)
	}
	if pMap["role"] != nil {
		role = pMap["role"].(string)
	}
	if pMap["email"] != nil {
		if emailStr, ok := pMap["email"].(string); ok {
			email = emailStr
		}
	}

	if email == "" && len(accountInfo.Emails) > 0 {
		for _, e := range accountInfo.Emails {
			if e.IsPrimary {
				email = e.Address
				break
			}
		}
		if email == "" {
			email = accountInfo.Emails[0].Address
		}
	}

	if email == "" {
		email = accountInfo.Login
	}

	return &client.CreateInvitationBody{
		Email:     email,
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
