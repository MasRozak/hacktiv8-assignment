package controllers

import (
	"net/http"

	"social-media-api/models"
	"social-media-api/services"

	"github.com/gin-gonic/gin"
)

var likeService = services.NewLikeService()

func CreateLike(c *gin.Context) {
	var like models.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := likeService.CreateLike(&like)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "post not found" {
			status = http.StatusNotFound
		} else if err.Error() == "satu user hanya boleh like satu post satu kali" {
			status = http.StatusConflict
		}

		c.JSON(status, models.Response{
			Message: "Failed to create like",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Message: "Like created successfully",
		Data:    like,
		Error:   nil,
	})
}

func GetLikesByPostID(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Post ID is required",
			Data:    nil,
			Error:   "missing post id",
		})
		return
	}

	likes, err := likeService.GetLikesByPostID(postID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch post likes",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Post likes retrieved successfully",
		Data:    likes,
		Error:   nil,
	})
}

func GetLikesByUserID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	likes, err := likeService.GetLikesByUserID(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch user likes",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "User likes retrieved successfully",
		Data:    likes,
		Error:   nil,
	})
}
