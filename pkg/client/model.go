package client

import "time"

type Response struct {
	Total int        `json:"total"`
	Limit int        `json:"limit"`
	Skip  int        `json:"skip"`
	Sys   SystemInfo `json:"sys"`
}

type GetUsersResponse struct {
	Response
	Items []User `json:"items"`
}

type User struct {
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	AvatarURL    string     `json:"avatarUrl"`
	Email        string     `json:"email"`
	SignupSource string     `json:"signupSource"`
	Activated    bool       `json:"activated"`
	SignInCount  int        `json:"signInCount"`
	Confirmed    bool       `json:"confirmed"`
	TwoFAEnabled bool       `json:"2faEnabled"`
	Sys          SystemInfo `json:"sys"`
}

type SystemInfo struct {
	Type      string    `json:"type"`
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy Link      `json:"createdBy"`
	UpdatedBy Link      `json:"updatedBy"`
	Org       Link      `json:"organization"`
	Team      Link      `json:"team"`
	Space     Link      `json:"space"`
	User      Link      `json:"user"`

	// these are only for org memberships
	LastActiveAt *time.Time  `json:"lastActiveAt"`
	Status       string      `json:"status"`
	SSO          interface{} `json:"sso"`

	// only for team memberships endpoint calls
	// https://www.contentful.com/developers/docs/references/user-management-api/#/reference/team-memberships
	OrganizationMembership Link `json:"organizationMembership"`
}

type Link struct {
	Sys LinkSys `json:"sys"`
}

type LinkSys struct {
	Type     string `json:"type"`
	LinkType string `json:"linkType"`
	ID       string `json:"id"`
}

type GetSpacesResponse struct {
	Response
	Items []Space `json:"items"`
}

type Space struct {
	Name string     `json:"name"`
	Sys  SystemInfo `json:"sys"`
}

type GetOrganizationsResponse struct {
	Response
	Items []Organization `json:"items"`
}

type Organization struct {
	Name string     `json:"name"`
	Sys  SystemInfo `json:"sys"`
}

type GetTeamsResponse struct {
	Response
	Items []Team `json:"items"`
}

type Team struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Sys         SystemInfo `json:"sys"`
}

type GetSpaceRolesResponse struct {
	Response
	Items []Role `json:"items"`
}

type Role struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Policies    []Policy    `json:"policies"`
	Permissions Permissions `json:"permissions"`
	Sys         SystemInfo  `json:"sys"`
}

type Policy struct {
	Effect     string         `json:"effect"`
	Actions    any            `json:"actions"` // Can be string "all" or []string
	Constraint map[string]any `json:"constraint"`
}

type Permissions struct {
	ContentModel    []string `json:"ContentModel"`
	Settings        []string `json:"Settings"`
	ContentDelivery any      `json:"ContentDelivery"` // Can be string "all" or []string
}

type GetOrganizationMembershipsResponse struct {
	Response
	Items []OrganizationMembership `json:"items"`
}

type OrganizationMembership struct {
	Role                       string     `json:"role"`
	IsExemptFromRestrictedMode bool       `json:"isExemptFromRestrictedMode"`
	Sys                        SystemInfo `json:"sys"`
}

type GetSpaceMembershipsResponse struct {
	Response
	Items []SpaceMembership `json:"items"`
}

type SpaceMembership struct {
	Admin bool       `json:"admin"`
	Sys   SystemInfo `json:"sys"`
	Roles []LinkRole `json:"roles"`
}

type LinkRole struct {
	Name string  `json:"name"`
	Sys  LinkSys `json:"sys"`
}

type GetTeamMembershipsResponse struct {
	Response
	Items []TeamMembership `json:"items"`
}

type TeamMembership struct {
	Sys SystemInfo `json:"sys"`
}
