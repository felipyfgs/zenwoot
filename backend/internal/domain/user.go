package domain

import "time"

type UserRole string

const (
	RoleAgent    UserRole = "agent"
	RoleAdmin    UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

type UserStatus string

const (
	UserActive    UserStatus = "active"
	UserInactive  UserStatus = "inactive"
	UserPending   UserStatus = "pending"
)

// User represents an agent/admin in the system.
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	DisplayName  string     `json:"displayName,omitempty"`
	AvatarURL    string     `json:"avatarUrl,omitempty"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	PasswordHash string    `json:"-"`
	Provider     string    `json:"provider,omitempty"` // sso, email, etc
	UID          string    `json:"uid,omitempty"`      // external provider ID
	Settings     any       `json:"settings"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// AccountUser links a user to an account with specific permissions.
type AccountUser struct {
	ID         string   `json:"id"`
	AccountID  string   `json:"accountId"`
	UserID     string   `json:"userId"`
	Role       UserRole `json:"role"` // can override user's global role per account
	Active     bool     `json:"active"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
