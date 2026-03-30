package dto

type UserCreateReq struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Role        string `json:"role"`
	Password    string `json:"password,omitempty"`
	Provider    string `json:"provider,omitempty"`
	UID         string `json:"uid,omitempty"`
}

type UserUpdateReq struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	Role        string `json:"role,omitempty"`
	Status      string `json:"status,omitempty"`
}

type TeamCreateReq struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	AllowAutoAssign bool   `json:"allowAutoAssign"`
}

type TeamUpdateReq struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	AllowAutoAssign *bool  `json:"allowAutoAssign,omitempty"`
}

type TeamMemberReq struct {
	Role string `json:"role,omitempty"`
}

type LabelCreateReq struct {
	Title       string `json:"title"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

type LabelUpdateReq struct {
	Title       string `json:"title,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

type SetLabelsReq struct {
	LabelIDs []string `json:"labelIds"`
}

type AssignReq struct {
	UserID string `json:"userId,omitempty"`
	TeamID string `json:"teamId,omitempty"`
}

type PriorityReq struct {
	Priority string `json:"priority"`
}

type SnoozeReq struct {
	Until string `json:"until"` // RFC3339 timestamp
}

type CannedResponseCreateReq struct {
	ShortCode string `json:"shortCode"`
	Content   string `json:"content"`
}

type CannedResponseUpdateReq struct {
	ShortCode string `json:"shortCode,omitempty"`
	Content   string `json:"content,omitempty"`
}

type NoteCreateReq struct {
	Content string `json:"content"`
}
