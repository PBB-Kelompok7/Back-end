package handler

import (
	"context"
	"crowdfunding-minpro-alterra/modules/user"
	"crowdfunding-minpro-alterra/utils/auth"
	"crowdfunding-minpro-alterra/utils/helper"
	"net/http"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
	cloudinary  *cloudinary.Cloudinary
}

func NewUserHandler(userService user.Service, authService auth.Service, cloudinary *cloudinary.Cloudinary) *userHandler {
	return &userHandler{userService, authService, cloudinary}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	newUser, err := h.userService.RegisterUser(input)
	
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered.", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	formatter := user.FormatUser(loggedinUser, token)

	response := helper.APIResponse("Login successfuly.", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server error."}

		response := helper.APIResponse("Email checking failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	data := gin.H{
		"is_available" : IsEmailAvailable,
	}

	metaMessage := "Email has been registered."

	if IsEmailAvailable {
		metaMessage = "Email is available."
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image.", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		response := helper.APIResponse("Failed to open uploaded file.", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer fileReader.Close()

	params := uploader.UploadParams{
		Folder:    "avatars",
		Overwrite: true,
	}

	uploadResult, err := h.cloudinary.Upload.Upload(context.Background(), fileReader, params)
	if err != nil {
		response := helper.APIResponse("Failed to upload avatar image.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	imageURL := uploadResult.SecureURL

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	currentUser.AvatarFileName = imageURL
	_, err = h.userService.SaveAvatar(userID, imageURL)
	if err != nil {
		response := helper.APIResponse("Failed to save avatar image URL.", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := user.GetFormatUser(currentUser)
	response := helper.APIResponse("Avatar uploaded successfully.", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	formatter := user.GetFormatUser(currentUser)

	response := helper.APIResponse("Successfuly fetch user.", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h * userHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()

	if err != nil {
		response := helper.APIResponse("Error to get all users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	if currentUser.ID != 2 {
		response := helper.APIResponse("You are not authorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	response := helper.APIResponse("List of all users", http.StatusOK, "success", user.GetFormatUsers(users))
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	var input struct {
			ID int `uri:"id" binding:"required"`
	}

	err := c.ShouldBindUri(&input)
	if err != nil {
			response := helper.APIResponse("Failed to delete user", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	if currentUser.ID != 2 {
			response := helper.APIResponse("You are not authorized", http.StatusForbidden, "error", nil)
			c.JSON(http.StatusForbidden, response)
			return
	}

	err = h.userService.DeleteUser(input.ID)
	if err != nil {
			response := helper.APIResponse("Failed to delete user", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	response := helper.APIResponse("User deleted successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
