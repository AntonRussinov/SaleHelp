package domain

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID          uuid.UUID  `json:"id"`
		CreatedAt   time.Time  `json:"createdAt"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		Name        string     `json:"name"`
		Email       string     `json:"email"`
		Roles       []UserRole `json:"-"`
		Role        UserRole   `json:"role"`
		InitialRole UserRole   `json:"-"`
		Active      bool       `json:"active"`
		Remark      string     `json:"remark"`
	}

	UserRole string
)

const (
	UserRoleAdmin      UserRole = "SandStorm.Admin"
	UserRoleUser       UserRole = "SandStorm.User"
	UserRoleViewer     UserRole = "SandStorm.Viewer"
	UserRoleEditor     UserRole = "SandStorm.Editor"
	UserRoleSuperAdmin UserRole = "SandStorm.SuperAdmin"
)

func rolesOrder() []UserRole {
	return []UserRole{
		UserRoleSuperAdmin,
		UserRoleAdmin,
		UserRoleEditor,
		UserRoleViewer,
		UserRoleUser,
	}
}

func GetMaxRole(roles []UserRole) UserRole {
	ordered := rolesOrder()

	for _, role := range ordered {
		if slices.Contains(roles, role) {
			return role
		}
	}

	return ""
}

func (u *User) GetMaxRole() UserRole {
	ordered := rolesOrder()

	for _, role := range ordered {
		if u.HasRole(role) {
			return role
		}
	}

	return ""
}

func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

func (u UserRole) IsValid() bool {
	return u == UserRoleAdmin ||
		u == UserRoleUser ||
		u == UserRoleEditor ||
		u == UserRoleViewer ||
		u == UserRoleSuperAdmin
}

func AllUserRoles() []UserRole {
	return []UserRole{
		UserRoleAdmin,
		UserRoleEditor,
		UserRoleSuperAdmin,
		UserRoleUser,
		UserRoleViewer,
	}
}
