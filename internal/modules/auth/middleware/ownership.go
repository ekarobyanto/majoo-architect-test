package middleware

import (
	"github.com/user/simple-blog/internal/modules/auth/domain"
)

// IsOwnerOrAdmin checks if the current user is the owner of a resource or an admin
func IsOwnerOrAdmin(user *domain.UserContext, resourceOwnerID string) bool {
	if user == nil {
		return false
	}

	// Check if user is admin
	for _, role := range user.Roles {
		if role == "admin" {
			return true
		}
	}

	// Check if user is owner
	return user.ID == resourceOwnerID
}
