package donation

import (
	"crowdfunding-minpro-alterra/modules/campaign"
	"crowdfunding-minpro-alterra/modules/payment"
	"errors"
	"strconv"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func (s *service) GetAllTransactions() ([]Donation, error) {
	panic("unimplemented")
}

type Service interface {
	GetDonationsByCampaignID(input GetCampaignDonationsInput) ([]Donation, error)
	GetDonationsByUserID(userID int) ([]Donation, error)
	CreateDonation(input CreateDonationInput) (Donation, error)
	ProcessPayment(input DonationNotificationInput) error
	GetAllTransactions() ([]Donation, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetDonationsByCampaignID(input GetCampaignDonationsInput) ([]Donation, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)

	if err != nil {
		return []Donation{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Donation{}, errors.New("Not an owner of the campaign.")
	}

	donations, err := s.repository.GetByCampaignID(input.ID)

	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (s *service) GetDonationsByUserID(userID int) ([]Donation, error) {
	donations, err := s.repository.GetByUserID(userID)

	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (s *service) CreateDonation(input CreateDonationInput) (Donation, error) {
	donation := Donation{}

	donation.CampaignID = input.CampaignID
	donation.Amount = input.Amount
	donation.UserID = input.User.ID
	donation.Status = "pending"
	// donation.Code = ""

	newDonation, err := s.repository.Save(donation)

	if err != nil {
		return newDonation, err
	}

	paymentDonation := payment.Donation{
		ID:     newDonation.ID,
		Amount: newDonation.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentDonation, input.User)

	if err != nil {
		return newDonation, err
	}

	newDonation.PaymentURL = paymentURL
	newDonation, err = s.repository.Update(newDonation)

	if err != nil {
		return newDonation, err
	}

	return newDonation, nil
}

func (s *service) ProcessPayment(input DonationNotificationInput) error {
	donation_id, _ := strconv.Atoi(input.OrderID)

	donation, err := s.repository.GetByID(donation_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		donation.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		donation.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		donation.Status = "cancelled"
	}

	updatedDonation, err := s.repository.Update(donation)

	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedDonation.CampaignID)

	if err != nil {
		return err
	}

	if updatedDonation.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedDonation.Amount

		_, err := s.campaignRepository.Update(campaign)

		if err != nil {
			return err
		}
	}

	return nil
}

// func (s *service) GetAllTransactions() ([]Donation, error) {
// 	donations, err := s.repository.FindAll()
// 	if err != nil {
// 		return donations, err
// 	}

// 	return donations, nil
// }
