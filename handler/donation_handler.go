package handler

import (
	"crowdfunding-minpro-alterra/modules/donation"
	"crowdfunding-minpro-alterra/modules/user"
	"crowdfunding-minpro-alterra/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type donationHandler struct {
	service donation.Service
}

func NewDonationHandler(service donation.Service) *donationHandler {
	return &donationHandler{service}
}

func (h *donationHandler) GetCampaignDonations(c *gin.Context) {
	var input donation.GetCampaignDonationsInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get campaign donations.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	donations, err := h.service.GetDonationsByCampaignID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get campaign donations.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	response := helper.APIResponse("Campaign donations.", http.StatusOK, "success", donation.FormatCampaignDonations(donations))
	c.JSON(http.StatusOK, response)
}

func (h *donationHandler) GetUserDonations(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	donations, err := h.service.GetDonationsByUserID(userID)

	if err != nil {
		response := helper.APIResponse("Failed to get user donations.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	response := helper.APIResponse("User donations.", http.StatusOK, "success", donation.FormatUserDonations(donations))
	c.JSON(http.StatusOK, response)
}

func (h *donationHandler) CreateDonation(c *gin.Context) {
	var logger = logrus.New()
	var input donation.CreateDonationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		logger.Error(err)
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create donation.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}
	logger.Info("input: ", input)

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newDonation, err := h.service.CreateDonation(input)

	if err != nil {
		logger.Error(err)
		response := helper.APIResponse("Failed to create donation.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	logger.Info("newDonation: ", newDonation)

	response := helper.APIResponse("Donation created.", http.StatusOK, "success", donation.FormatDonation(newDonation))
	c.JSON(http.StatusOK, response)
}

func (h *donationHandler) GetNotification(c *gin.Context) {
	var input donation.DonationNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get notification.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	err = h.service.ProcessPayment(input)

	if err != nil {
		response := helper.APIResponse("Failed to get notification.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	c.JSON(http.StatusOK, input)
}
