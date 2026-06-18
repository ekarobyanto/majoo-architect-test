package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	_ "github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/internal/platform/response"
	"github.com/user/simple-blog/internal/platform/validation"
)

type PostHandler struct {
	svc domain.PostService
}

func NewPostHandler(svc domain.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// Create godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreatePostRequest true "Post details"
// @Success 201 {object} response.Response{data=models.Post}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Router /posts [post]
func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req domain.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	user := c.Locals("user").(*authDomain.UserContext)
	post, err := h.svc.Create(c.Context(), user.ID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, "Post created successfully", post)
}

// GetByID godoc
// @Summary Get post by ID
// @Description Get a post by its ID
// @Tags Posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response{data=models.Post}
// @Failure 404 {object} errors.ErrorResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, "Post retrieved", post)
}

// GetPaginated godoc
// @Summary Get all posts
// @Description Get paginated list of posts
// @Tags Posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=domain.PaginatedPostResponse}
// @Router /posts [get]
func (h *PostHandler) GetPaginated(c *fiber.Ctx) error {
	var query domain.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid query parameters")
	}

	if err := validation.Validate(query); err != nil {
		return err
	}

	resp, err := h.svc.GetPaginated(c.Context(), query)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Posts retrieved", resp)
}

// Update godoc
// @Summary Update a post
// @Description Update an existing post by ID
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body domain.UpdatePostRequest true "Update details"
// @Success 200 {object} response.Response{data=models.Post}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	user := c.Locals("user").(*authDomain.UserContext)
	post, err := h.svc.Update(c.Context(), id, user, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Post updated successfully", post)
}

// Delete godoc
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*authDomain.UserContext)

	if err := h.svc.Delete(c.Context(), id, user); err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Post deleted successfully", nil)
}
