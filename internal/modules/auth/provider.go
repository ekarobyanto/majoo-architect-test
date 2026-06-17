package auth

import (
	"github.com/google/wire"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/handler"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/repository"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/service"
)

var ProviderSet = wire.NewSet(
	repository.NewAuthRepository,
	service.NewAuthService,
	handler.NewAuthHandler,
)
