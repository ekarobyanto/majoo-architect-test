package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/platform/response"
	"github.com/user/simple-blog/internal/platform/validation"
)

type CommentHandler struct {
	svc domain.CommentService
}

func NewCommentHandler(svc domain.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// Create godoc
// @Summary Add a comment to a post
// @Description Add a new comment to an existing blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Post ID"
// @Param request body domain.CreateCommentRequest true "Comment details"
// @Success 201 {object} response.Response{data=models.Comment}
// @Router /posts/{id}/comments [post]
func (h *CommentHandler) Create(c *fiber.Ctx) error {
	postID := c.Params("id")
	var req domain.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	user := c.Locals("user").(*authDomain.UserContext)
	comment, err := h.svc.Create(c.Context(), postID, user.ID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, "Comment added successfully", comment)
}

// Update godoc
// @Summary Update a comment
// @Description Update an existing comment by ID
// @Tags Comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Comment ID"
// @Param request body domain.UpdateCommentRequest true "Update details"
// @Success 200 {object} response.Response{data=models.Comment}
// @Router /comments/{id} [put]
func (h *CommentHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	user := c.Locals("user").(*authDomain.UserContext)
	comment, err := h.svc.Update(c.Context(), id, user, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Comment updated successfully", comment)
}

// Delete godoc
// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags Comments
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Comment ID"
// @Success 200 {object} response.Response
// @Router /comments/{id} [delete]
func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*authDomain.UserContext)

	if err := h.svc.Delete(c.Context(), id, user); err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Comment deleted successfully", nil)
}
