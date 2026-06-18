package authorization_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/simple-blog/internal/modules/auth/authorization"
	"github.com/user/simple-blog/internal/modules/auth/domain"
)

func TestIsOwnerOrAdmin(t *testing.T) {
	adminUser := &domain.UserContext{ID: "admin-1", Roles: []string{"admin"}}
	regularUser := &domain.UserContext{ID: "user-1", Roles: []string{"user"}}
	anotherUser := &domain.UserContext{ID: "user-2", Roles: []string{"user"}}

	resourceOwnerID := "user-1"

	assert.True(t, authorization.IsOwnerOrAdmin(adminUser, resourceOwnerID), "Admin should have access")
	assert.True(t, authorization.IsOwnerOrAdmin(regularUser, resourceOwnerID), "Owner should have access")
	assert.False(t, authorization.IsOwnerOrAdmin(anotherUser, resourceOwnerID), "Non-owner should not have access")
	assert.False(t, authorization.IsOwnerOrAdmin(nil, resourceOwnerID), "Nil user should not have access")
}
