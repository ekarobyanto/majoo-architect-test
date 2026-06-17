package health

import (
	"github.com/google/wire"
	"github.com/user/go-backend-boilerplate/internal/modules/health/handler"
	"github.com/user/go-backend-boilerplate/internal/modules/health/repository"
	"github.com/user/go-backend-boilerplate/internal/modules/health/service"
)

var ProviderSet = wire.NewSet(
	repository.NewHealthRepository,
	service.NewHealthService,
	handler.NewHealthHandler,
)
