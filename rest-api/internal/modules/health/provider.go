package health

import (
	"github.com/google/wire"
	"github.com/user/simple-blog/internal/modules/health/handler"
	"github.com/user/simple-blog/internal/modules/health/repository"
	"github.com/user/simple-blog/internal/modules/health/service"
)

var ProviderSet = wire.NewSet(
	repository.NewHealthRepository,
	service.NewHealthService,
	handler.NewHealthHandler,
)
