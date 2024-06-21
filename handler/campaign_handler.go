package handler

import (
	"context"
	"crowdfunding-minpro-alterra/modules/campaign"
	"crowdfunding-minpro-alterra/modules/user"
	"crowdfunding-minpro-alterra/utils/helper"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
	cloudinary  *cloudinary.Cloudinary
}

func NewCampaignHandler(service campaign.Service, cloudinary *cloudinary.Cloudinary) *campaignHandler {
	return &campaignHandler{service, cloudinary}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	response := helper.APIResponse("List of campaigns.", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign.", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)

		return
	}

	response := helper.APIResponse("Campaign detail.", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign.", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
			response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
	}


	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to upload campaign image.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image.", http.StatusBadRequest, "error", data)
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
		Folder:    "campaigns",
		Overwrite: true,
	}

	uploadResult, err := h.cloudinary.Upload.Upload(context.Background(), fileReader, params)
		if err != nil {
		response := helper.APIResponse("Failed to upload campaign image.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	imageURL := uploadResult.SecureURL

	_, err = h.service.SaveCampaignImage(input, imageURL)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image.", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image uploaded successfully.", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
