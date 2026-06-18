package authorization

import "github.com/user/simple-blog/internal/modules/auth/domain"

// IsOwnerOrAdmin checks if the current user owns a resource or has the admin role.
func IsOwnerOrAdmin(user *domain.UserContext, resourceOwnerID string) bool {
	if user == nil {
		return false
	}

	for _, role := range user.Roles {
		if role == "admin" {
			return true
		}
	}

	return user.ID == resourceOwnerID
}
