package donation

import (
	"crowdfunding-minpro-alterra/modules/campaign"
	"crowdfunding-minpro-alterra/modules/user"

	"time"

	"github.com/leekchan/accounting"
)

type Donation struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User
	Campaign   campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt *time.Time
}

func (d Donation) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(d.Amount)
}