package donation

import "crowdfunding-minpro-alterra/modules/user"

type GetCampaignDonationsInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}

type CreateDonationInput struct {
	Amount int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User user.User
}

type DonationNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}