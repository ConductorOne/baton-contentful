package client

import "time"

type Response struct {
	Total int        `json:"total"`
	Limit int        `json:"limit"`
	Skip  int        `json:"skip"`
	Sys   SystemInfo `json:"sys"`
}

// Response represents the top-level JSON structure
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
