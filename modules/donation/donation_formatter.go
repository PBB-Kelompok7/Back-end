package donation

import "time"

type CampaignDonationFormatter struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignDonation(donation Donation) CampaignDonationFormatter {
	formatter := CampaignDonationFormatter{}

	formatter.ID = donation.ID
	formatter.Name = donation.User.Name
	formatter.Amount = donation.Amount
	formatter.CreatedAt = donation.CreatedAt

	return formatter
}

func FormatCampaignDonations(donations []Donation) []CampaignDonationFormatter {

	if len(donations) == 0 {
		return []CampaignDonationFormatter{}
	}

	var donationsFormatter []CampaignDonationFormatter

	for _, donation := range donations {
		formatter := FormatCampaignDonation(donation)
		donationsFormatter = append(donationsFormatter, formatter)
	}

	return donationsFormatter
}

type UserDonationFormatter struct {
	ID        int    `json:"id"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserDonation(donation Donation) UserDonationFormatter {
	formatter := UserDonationFormatter{}

	formatter.ID = donation.ID
	formatter.Amount = donation.Amount
	formatter.Status = donation.Status
	formatter.CreatedAt = donation.CreatedAt

	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Name = donation.Campaign.Name
	campaignFormatter.ImageURL = ""

	if len(donation.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = donation.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = campaignFormatter

	return formatter
}

func FormatUserDonations(donations []Donation) []UserDonationFormatter {

	if len(donations) == 0 {
		return []UserDonationFormatter{}
	}

	var donationsFormatter []UserDonationFormatter

	for _, donation := range donations {
		formatter := FormatUserDonation(donation)
		donationsFormatter = append(donationsFormatter, formatter)
	}

	return donationsFormatter
}

type DonationFormatter struct {
	ID        int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserID    int    `json:"user_id"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	Code      string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatDonation(donation Donation) DonationFormatter {
	formatter := DonationFormatter{}

	formatter.ID = donation.ID
	formatter.CampaignID = donation.CampaignID
	formatter.UserID = donation.UserID
	formatter.Amount = donation.Amount
	formatter.Status = donation.Status
	formatter.Code = donation.Code
	formatter.PaymentURL = donation.PaymentURL
	
	return formatter
}