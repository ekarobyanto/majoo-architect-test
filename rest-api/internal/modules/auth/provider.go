package auth

import (
	"github.com/google/wire"
	"github.com/user/simple-blog/internal/modules/auth/handler"
	"github.com/user/simple-blog/internal/modules/auth/repository"
	"github.com/user/simple-blog/internal/modules/auth/service"
)

var ProviderSet = wire.NewSet(
	repository.NewAuthRepository,
	service.NewAuthService,
	handler.NewAuthHandler,
)
