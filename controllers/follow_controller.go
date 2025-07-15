package controllers

import (
	"net/http"

	"social-media-api/models"
	"social-media-api/services"

	"github.com/gin-gonic/gin"
)

var followService = services.NewFollowService()

func CreateFollow(c *gin.Context) {
	var follow models.Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := followService.CreateFollow(&follow)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "cannot follow yourself" {
			status = http.StatusBadRequest
		} else if err.Error() == "follower user not found" || err.Error() == "following user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "already following this user" {
			status = http.StatusConflict
		}

		c.JSON(status, models.Response{
			Message: "Failed to create follow",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Message: "Follow created successfully",
		Data:    follow,
		Error:   nil,
	})
}

func DeleteFollow(c *gin.Context) {
	var follow models.Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := followService.DeleteFollow(follow.FollowerID, follow.FollowingID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "follow relationship not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to delete follow",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Follow deleted successfully",
		Data:    nil,
		Error:   nil,
	})
}

func GetFollowers(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	followers, err := followService.GetFollowers(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch followers",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Followers retrieved successfully",
		Data:    followers,
		Error:   nil,
	})
}

func GetFollowing(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	following, err := followService.GetFollowing(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch following",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Following retrieved successfully",
		Data:    following,
		Error:   nil,
	})
}
