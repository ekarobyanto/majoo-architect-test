package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthService_Check(t *testing.T) {
	repo := new(mockRepository)
	svc := NewHealthService(repo)

	t.Run("success", func(t *testing.T) {
		repo.On("Ping", mock.Anything).Return(nil).Once()
		resp, err := svc.Check(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "UP", resp.Status)
		assert.Equal(t, "Database connection is healthy", resp.Message)
	})

	t.Run("error", func(t *testing.T) {
		repo.On("Ping", mock.Anything).Return(errors.New("db error")).Once()
		resp, err := svc.Check(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "DOWN", resp.Status)
		assert.Equal(t, "Database connection failed", resp.Message)
	})
}
