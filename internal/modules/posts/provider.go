package posts

import (
	"github.com/google/wire"
	"github.com/user/simple-blog/internal/modules/posts/handler"
	"github.com/user/simple-blog/internal/modules/posts/repository"
	"github.com/user/simple-blog/internal/modules/posts/service"
)

var ProviderSet = wire.NewSet(
	repository.NewPostRepository,
	service.NewPostService,
	handler.NewPostHandler,
)
