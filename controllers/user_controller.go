package controllers

import (
	"net/http"

	"social-media-api/models"
	"social-media-api/services"

	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService()

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := userService.CreateUser(&user)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "username or email already exists" {
			status = http.StatusConflict
		} else if err.Error() == "username is required" || err.Error() == "email is required" || err.Error() == "invalid email format" {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.Response{
			Message: "Failed to create user",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Message: "User created successfully",
		Data:    user,
		Error:   nil,
	})
}

func GetAllUsers(c *gin.Context) {
	users, err := userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to fetch users",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "Users retrieved successfully",
		Data:    users,
		Error:   nil,
	})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	user, err := userService.GetUserByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to fetch user",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "User retrieved successfully",
		Data:    user,
		Error:   nil,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid JSON format",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	err := userService.UpdateUser(id, &user)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "username or email already exists" {
			status = http.StatusConflict
		} else if err.Error() == "username is required" || err.Error() == "email is required" || err.Error() == "invalid email format" {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.Response{
			Message: "Failed to update user",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "User updated successfully",
		Data:    user,
		Error:   nil,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Message: "User ID is required",
			Data:    nil,
			Error:   "missing user id",
		})
		return
	}

	err := userService.DeleteUser(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.Response{
			Message: "Failed to delete user",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "User deleted successfully",
		Data:    nil,
		Error:   nil,
	})
}
