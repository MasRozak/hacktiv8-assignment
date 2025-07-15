package controllers

import (
	"net/http"

	"social-media-api/models"
	"social-media-api/services"

	"github.com/gin-gonic/gin"
)

var postService = services.NewPostService()

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := postService.CreatePost(&post)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "content tidak boleh kosong" {
			status = http.StatusBadRequest
		} else if err.Error() == "user harus valid" {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.Response{
			Message: "Failed to create post",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Message: "Post created successfully",
		Data:    post,
		Error:   nil,
	})
}

func GetAllPosts(c *gin.Context) {
	userID := c.Query("user_id")
	keyword := c.Query("keyword")

	var posts []models.Post
	var err error

	if userID != "" || keyword != "" {
		posts, err = postService.GetPostsWithFilters(userID, keyword)
	} else {
		posts, err = postService.GetAllPosts()
	}

	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch posts",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Posts retrieved successfully",
		Data:    posts,
		Error:   nil,
	})
}

func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Post ID is required",
			Data:    nil,
			Error:   "missing post id",
		})
		return
	}

	post, err := postService.GetPostByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch post",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Post retrieved successfully",
		Data:    post,
		Error:   nil,
	})
}

func GetPostsByUserID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	posts, err := postService.GetPostsByUserID(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch user posts",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "User posts retrieved successfully",
		Data:    posts,
		Error:   nil,
	})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Post ID is required",
			Data:    nil,
			Error:   "missing post id",
		})
		return
	}

	err := postService.DeletePost(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to delete post",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Post deleted successfully",
		Data:    nil,
		Error:   nil,
	})
}
