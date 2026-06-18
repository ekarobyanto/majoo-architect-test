package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/middleware"
)

func TestIsOwnerOrAdmin(t *testing.T) {
	adminUser := &domain.UserContext{ID: "admin-1", Roles: []string{"admin"}}
	regularUser := &domain.UserContext{ID: "user-1", Roles: []string{"user"}}
	anotherUser := &domain.UserContext{ID: "user-2", Roles: []string{"user"}}

	resourceOwnerID := "user-1"

	assert.True(t, middleware.IsOwnerOrAdmin(adminUser, resourceOwnerID), "Admin should have access")
	assert.True(t, middleware.IsOwnerOrAdmin(regularUser, resourceOwnerID), "Owner should have access")
	assert.False(t, middleware.IsOwnerOrAdmin(anotherUser, resourceOwnerID), "Non-owner should not have access")
	assert.False(t, middleware.IsOwnerOrAdmin(nil, resourceOwnerID), "Nil user should not have access")
}
