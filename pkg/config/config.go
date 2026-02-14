package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	TokenField = field.StringField(
		"token",
		field.WithDisplayName("Contentful API token"),
		field.WithDescription("The API token used to authenticate with the service."),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)

	OrgIdField = field.StringField(
		"organization-id",
		field.WithDisplayName("Organization ID"),
		field.WithDescription("The ID of the organization to use."),
		field.WithRequired(true),
	)

	BaseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Contentful API URL (for testing)"),
		field.WithHidden(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run. Note: these fields can be marked as optional or
	// required.
	ConfigurationFields = []field.SchemaField{
		TokenField,
		OrgIdField,
		BaseURLField,
	}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated. For example, a
	// username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Contentful"),
	field.WithHelpUrl("/docs/baton/contentful"),
	field.WithIconUrl("/static/app-icons/contentful.svg"),
)
