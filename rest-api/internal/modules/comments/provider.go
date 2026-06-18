package comments

import (
	"github.com/google/wire"
	"github.com/user/simple-blog/internal/modules/comments/handler"
	"github.com/user/simple-blog/internal/modules/comments/repository"
	"github.com/user/simple-blog/internal/modules/comments/service"
)

var ProviderSet = wire.NewSet(
	repository.NewCommentRepository,
	service.NewCommentService,
	handler.NewCommentHandler,
)
