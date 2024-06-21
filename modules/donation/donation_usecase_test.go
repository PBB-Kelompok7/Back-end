package donation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	GetByCampaignIDFunc func(CampaignID int) ([]Donation, error)
	GetByUserIDFunc     func(UserID int) ([]Donation, error)
	GetByIDFunc         func(ID int) (Donation, error)
	SaveFunc            func(donation Donation) (Donation, error)
	UpdateFunc          func(donation Donation) (Donation, error)
}

func (m *MockRepository) GetByCampaignID(CampaignID int) ([]Donation, error) {
	if m.GetByCampaignIDFunc != nil {
		return m.GetByCampaignIDFunc(CampaignID)
	}
	return nil, nil
}

func (m *MockRepository) GetByUserID(UserID int) ([]Donation, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(UserID)
	}
	return nil, nil
}

func (m *MockRepository) GetByID(ID int) (Donation, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ID)
	}
	return Donation{}, nil
}

func (m *MockRepository) Save(donation Donation) (Donation, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(donation)
	}
	return Donation{}, nil
}

func (m *MockRepository) Update(donation Donation) (Donation, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(donation)
	}
	return Donation{}, nil
}

func TestService_GetDonationsByUserID(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo, nil, nil)

	t.Run("Test GetDonationsByUserID with valid user ID", func(t *testing.T) {
		mockUserID := 1
		mockDonations := []Donation{
			{ID: 1, UserID: mockUserID},
			{ID: 2, UserID: mockUserID},
		}

		repo.GetByUserIDFunc = func(UserID int) ([]Donation, error) {
			if UserID == mockUserID {
				return mockDonations, nil
			}
			return nil, errors.New("user not found")
		}

		donations, err := service.GetDonationsByUserID(mockUserID)

		assert.NoError(t, err)
		assert.Equal(t, len(mockDonations), len(donations))
	})

	t.Run("Test GetDonationsByUserID with invalid user ID", func(t *testing.T) {
		repo.GetByUserIDFunc = func(UserID int) ([]Donation, error) {
			return nil, errors.New("user not found")
		}

		donations, err := service.GetDonationsByUserID(999)

		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
		assert.Empty(t, donations)
	})
}

