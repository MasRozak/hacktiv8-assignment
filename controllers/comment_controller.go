package controllers

import (
	"net/http"

	"social-media-api/models"
	"social-media-api/services"

	"github.com/gin-gonic/gin"
)

var commentService = services.NewCommentService()

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := commentService.CreateComment(&comment)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "content tidak boleh kosong" {
			status = http.StatusBadRequest
		} else if err.Error() == "user not found" || err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to create comment",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Message: "Comment created successfully",
		Data:    comment,
		Error:   nil,
	})
}

func GetCommentsByPostID(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Post ID is required",
			Data:    nil,
			Error:   "missing post id",
		})
		return
	}

	comments, err := commentService.GetCommentsByPostID(postID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch comments",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Comments retrieved successfully",
		Data:    comments,
		Error:   nil,
	})
}
