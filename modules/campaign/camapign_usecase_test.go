package campaign

import (
	"crowdfunding-minpro-alterra/modules/user"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	FindAllFunc               func() ([]Campaign, error)
	FindByUserIDFunc          func(userID int) ([]Campaign, error)
	FindByIDFunc              func(ID int) (Campaign, error)
	SaveFunc                  func(campaign Campaign) (Campaign, error)
	UpdateFunc                func(campaign Campaign) (Campaign, error)
	CreateImageFunc           func(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllImagesAsNonPrimaryFunc func(campaignID int) (bool, error)
}

func (m *MockRepository) FindAll() ([]Campaign, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

func (m *MockRepository) FindByUserID(userID int) ([]Campaign, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *MockRepository) FindByID(ID int) (Campaign, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ID)
	}
	return Campaign{}, nil
}

func (m *MockRepository) Save(campaign Campaign) (Campaign, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(campaign)
	}
	return Campaign{}, nil
}

func (m *MockRepository) Update(campaign Campaign) (Campaign, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(campaign)
	}
	return Campaign{}, nil
}

func (m *MockRepository) CreateImage(campaignImage CampaignImage) (CampaignImage, error) {
	if m.CreateImageFunc != nil {
		return m.CreateImageFunc(campaignImage)
	}
	return CampaignImage{}, nil
}

func (m *MockRepository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	if m.MarkAllImagesAsNonPrimaryFunc != nil {
		return m.MarkAllImagesAsNonPrimaryFunc(campaignID)
	}
	return false, nil
}

func TestGetCampaigns(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	t.Run("Test GetCampaigns for specific user", func(t *testing.T) {
		mockUserID := 1

		mockCampaigns := []Campaign{
			{ID: 1, Name: "Campaign 1", UserID: mockUserID},
			{ID: 2, Name: "Campaign 2", UserID: mockUserID},
		}

		repo.FindByUserIDFunc = func(userID int) ([]Campaign, error) {
			if userID == mockUserID {
				return mockCampaigns, nil
			}
			return []Campaign{}, errors.New("user not found")
		}

		campaigns, err := service.GetCampaigns(mockUserID)

		assert.NoError(t, err)
		assert.Equal(t, len(mockCampaigns), len(campaigns))
		for i := range campaigns {
			assert.Equal(t, mockCampaigns[i].ID, campaigns[i].ID)
			assert.Equal(t, mockCampaigns[i].Name, campaigns[i].Name)
			assert.Equal(t, mockCampaigns[i].UserID, campaigns[i].UserID)
		}
	})

	t.Run("Test GetCampaigns for all users", func(t *testing.T) {
		mockCampaigns := []Campaign{
			{ID: 1, Name: "Campaign 1", UserID: 1},
			{ID: 2, Name: "Campaign 2", UserID: 2},
		}

		repo.FindAllFunc = func() ([]Campaign, error) {
			return mockCampaigns, nil
		}

		campaigns, err := service.GetCampaigns(0)

		assert.NoError(t, err)
		assert.Equal(t, len(mockCampaigns), len(campaigns))
		for i := range campaigns {
			assert.Equal(t, mockCampaigns[i].ID, campaigns[i].ID)
			assert.Equal(t, mockCampaigns[i].Name, campaigns[i].Name)
			assert.Equal(t, mockCampaigns[i].UserID, campaigns[i].UserID)
		}
	})

	t.Run("Test GetCampaigns with invalid user ID", func(t *testing.T) {
		repo.FindByUserIDFunc = func(userID int) ([]Campaign, error) {
			return []Campaign{}, errors.New("user not found")
		}

		campaigns, err := service.GetCampaigns(999)

		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
		assert.Equal(t, []Campaign{}, campaigns)
	})

	t.Run("Test GetCampaigns with repository error", func(t *testing.T) {
		repo.FindAllFunc = func() ([]Campaign, error) {
			return []Campaign{}, errors.New("repository error")
		}

		campaigns, err := service.GetCampaigns(0)

		assert.Error(t, err)
		assert.EqualError(t, err, "repository error")
		assert.Equal(t, []Campaign{}, campaigns)
	})
}

func TestGetCampaignByID(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	t.Run("Test GetCampaignByID for existing campaign", func(t *testing.T) {
		mockCampaignID := 1

		mockCampaign := Campaign{
			ID:          mockCampaignID,
			Name:        "Campaign 1",
			Description: "Description of Campaign 1",
			UserID:      1,
		}

		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			if ID == mockCampaignID {
				return mockCampaign, nil
			}
			return Campaign{}, errors.New("campaign not found")
		}

		campaign, err := service.GetCampaignByID(GetCampaignDetailInput{ID: mockCampaignID})

		assert.NoError(t, err)
		assert.Equal(t, mockCampaignID, campaign.ID)
		assert.Equal(t, mockCampaign.Name, campaign.Name)
		assert.Equal(t, mockCampaign.Description, campaign.Description)
		assert.Equal(t, mockCampaign.UserID, campaign.UserID)
	})

	t.Run("Test GetCampaignByID for non-existing campaign", func(t *testing.T) {
		mockCampaignID := 999

		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			return Campaign{}, errors.New("campaign not found")
		}

		campaign, err := service.GetCampaignByID(GetCampaignDetailInput{ID: mockCampaignID})

		assert.Error(t, err)
		assert.EqualError(t, err, "campaign not found")
		assert.Equal(t, Campaign{}, campaign)
	})

	t.Run("Test GetCampaignByID with repository error", func(t *testing.T) {
		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			return Campaign{}, errors.New("repository error")
		}

		campaign, err := service.GetCampaignByID(GetCampaignDetailInput{ID: 1})

		assert.Error(t, err)
		assert.EqualError(t, err, "repository error")
		assert.Equal(t, Campaign{}, campaign)
	})
}

func TestCreateCampaign(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	t.Run("Test CreateCampaign success", func(t *testing.T) {
		mockInput := CreateCampaignInput{
			Name:             "Campaign 1",
			ShortDescription: "Short description",
			Description:      "Description",
			GoalAmount:       1000,
			User:             user.User{ID: 1},
		}

		expectedCampaign := Campaign{
			Name:             mockInput.Name,
			ShortDescription: mockInput.ShortDescription,
			Description:      mockInput.Description,
			GoalAmount:       mockInput.GoalAmount,
			UserID:           mockInput.User.ID,
		}

		repo.SaveFunc = func(campaign Campaign) (Campaign, error) {
			return expectedCampaign, nil
		}

		newCampaign, err := service.CreateCampaign(mockInput)

		assert.NoError(t, err)
		assert.NotNil(t, newCampaign)
		assert.Equal(t, expectedCampaign, newCampaign)
	})

	t.Run("Test CreateCampaign with error", func(t *testing.T) {
		repo.SaveFunc = func(campaign Campaign) (Campaign, error) {
			return Campaign{}, errors.New("unable to save campaign")
		}

		newCampaign, err := service.CreateCampaign(CreateCampaignInput{})

		assert.Error(t, err)
		assert.EqualError(t, err, "unable to save campaign")
		assert.Equal(t, Campaign{}, newCampaign)
	})
}

func TestUpdateCampaign(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	t.Run("Test UpdateCampaign success", func(t *testing.T) {
		mockInputID := 1
		mockInputData := CreateCampaignInput{
			Name:             "Updated Campaign",
			ShortDescription: "Updated short description",
			Description:      "Updated description",
			GoalAmount:       2000,
			User:             user.User{ID: 1},
		}

		existingCampaign := Campaign{
			ID:               mockInputID,
			Name:             "Campaign 1",
			ShortDescription: "Short description",
			Description:      "Description",
			GoalAmount:       1000,
			UserID:           1,
		}

		expectedUpdatedCampaign := Campaign{
			ID:               mockInputID,
			Name:             mockInputData.Name,
			ShortDescription: mockInputData.ShortDescription,
			Description:      mockInputData.Description,
			GoalAmount:       mockInputData.GoalAmount,
			UserID:           mockInputData.User.ID,
		}

		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			if ID == mockInputID {
				return existingCampaign, nil
			}
			return Campaign{}, errors.New("campaign not found")
		}

		repo.UpdateFunc = func(campaign Campaign) (Campaign, error) {
			return expectedUpdatedCampaign, nil
		}

		updatedCampaign, err := service.UpdateCampaign(GetCampaignDetailInput{ID: mockInputID}, mockInputData)

		assert.NoError(t, err)
		assert.NotNil(t, updatedCampaign)
		assert.Equal(t, expectedUpdatedCampaign, updatedCampaign)
	})

	t.Run("Test UpdateCampaign with non-existing ID", func(t *testing.T) {
		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			return Campaign{}, errors.New("campaign not found")
		}

		updatedCampaign, err := service.UpdateCampaign(GetCampaignDetailInput{ID: 999}, CreateCampaignInput{})

		assert.Error(t, err)
		assert.EqualError(t, err, "campaign not found")
		assert.Equal(t, Campaign{}, updatedCampaign)
	})
}

func TestUpdateCampaign_NotOwner(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	mockCampaignID := 1
	mockUserID := 2 
	inputID := GetCampaignDetailInput{ID: mockCampaignID}
	inputData := CreateCampaignInput{Name: "Updated Campaign", User: user.User{ID: mockUserID}}

	repo.FindByIDFunc = func(ID int) (Campaign, error) {
		return Campaign{ID: mockCampaignID, UserID: 1}, nil // Assume existing campaign owned by user with ID 1
	}

	_, err := service.UpdateCampaign(inputID, inputData)

	assert.Error(t, err)
	assert.EqualError(t, err, "Not an owner of the campaign.")
}

func TestSaveCampaignImage(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	t.Run("Test SaveCampaignImage success", func(t *testing.T) {
		mockCampaignID := 1
		mockFileLocation := "/path/to/image.jpg"
		mockIsPrimary := true
		mockUser := user.User{ID: 1}

		existingCampaign := Campaign{ID: mockCampaignID, UserID: mockUser.ID}

		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			if ID == mockCampaignID {
				return existingCampaign, nil
			}
			return Campaign{}, errors.New("campaign not found")
		}

		repo.CreateImageFunc = func(campaignImage CampaignImage) (CampaignImage, error) {
			var isPrimaryInt int
			if mockIsPrimary {
				isPrimaryInt = 1
			} else {
				isPrimaryInt = 0
			}
			return CampaignImage{
				CampaignID: mockCampaignID,
				FileName:   mockFileLocation,
				IsPrimary:  isPrimaryInt, 
			}, nil
		}

		var isPrimaryInt int

		newCampaignImage, err := service.SaveCampaignImage(CreateCampaignImageInput{
			CampaignID: mockCampaignID,
			IsPrimary:  mockIsPrimary,
			User:       mockUser,
		}, mockFileLocation)

		if mockIsPrimary {
			isPrimaryInt = 1
		} else {
			isPrimaryInt = 0
		}

		assert.NoError(t, err)
		assert.NotNil(t, newCampaignImage)
		assert.Equal(t, mockFileLocation, newCampaignImage.FileName)
		assert.Equal(t, isPrimaryInt, newCampaignImage.IsPrimary) // Check if isPrimaryInt is equal to newCampaignImage.IsPrimary
	})

	t.Run("Test SaveCampaignImage with non-existing campaign", func(t *testing.T) {
		repo.FindByIDFunc = func(ID int) (Campaign, error) {
			return Campaign{}, errors.New("campaign not found")
		}

		newCampaignImage, err := service.SaveCampaignImage(CreateCampaignImageInput{}, "/path/to/image.jpg")

		assert.Error(t, err)
		assert.EqualError(t, err, "campaign not found")
		assert.Equal(t, CampaignImage{}, newCampaignImage)
	})
}

func TestSaveCampaignImage_NotOwner(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	mockUser := user.User{
		ID:   1,
		Name: "John",
	}

	mockCampaignID := 1
	mockCampaign := Campaign{
		ID:     mockCampaignID,
		UserID: 2, // Non-owner user ID
	}

	repo.FindByIDFunc = func(ID int) (Campaign, error) {
		if ID == mockCampaignID {
			return mockCampaign, nil
		}
		return Campaign{}, errors.New("campaign not found")
	}

	input := CreateCampaignImageInput{
		CampaignID: mockCampaignID,
		User:       mockUser,
	}

	_, err := service.SaveCampaignImage(input, "file_location.jpg")

	assert.Error(t, err)
	assert.EqualError(t, err, "Not an owner of the campaign.")
}

