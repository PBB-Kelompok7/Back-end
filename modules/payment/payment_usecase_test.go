package payment

import (
	"crowdfunding-minpro-alterra/modules/user"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veritrans/go-midtrans"
)

type MockSnapGateway struct {
	Client *midtrans.Client
}

func (m *MockSnapGateway) GetToken(req *midtrans.SnapReq) (*midtrans.SnapResponse, error) {
	return &midtrans.SnapResponse{
		RedirectURL: "http://example.com/payment",
	}, nil
}

type MockService struct {
	GetTokenFunc func(req *midtrans.SnapReq) (*midtrans.SnapResponse, error)
}

func (m *MockService) GetPaymentURL(donation Donation, user user.User) (string, error) {
	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(donation.ID),
			GrossAmt: int64(donation.Amount),
		},
	}

	resp, err := m.GetTokenFunc(snapReq)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil
}

func TestService_GetPaymentURL(t *testing.T) {
	os.Setenv("MIDTRANS_SERVER_KEY", "server_key")
	os.Setenv("MIDTRANS_CLIENT_KEY", "client_key")

	mockDonation := Donation{
		ID:     123,
		Amount: 1000,
	}

	mockUser := user.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockService := &MockService{
		GetTokenFunc: func(req *midtrans.SnapReq) (*midtrans.SnapResponse, error) {
			assert.Equal(t, mockUser.Email, req.CustomerDetail.Email)
			assert.Equal(t, mockUser.Name, req.CustomerDetail.FName)
			assert.Equal(t, strconv.Itoa(mockDonation.ID), req.TransactionDetails.OrderID)
			assert.Equal(t, int64(mockDonation.Amount), req.TransactionDetails.GrossAmt)

			return &midtrans.SnapResponse{
				RedirectURL: "http://example.com/payment",
			}, nil
		},
	}

	url, err := mockService.GetPaymentURL(mockDonation, mockUser)

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/payment", url)
}

func TestMockService_GetToken(t *testing.T) {
	mockReq := &midtrans.SnapReq{}
	mockResp := &midtrans.SnapResponse{
		RedirectURL: "http://example.com/payment",
	}

	mockService := &MockService{
		GetTokenFunc: func(req *midtrans.SnapReq) (*midtrans.SnapResponse, error) {
			assert.Equal(t, mockReq, req)
			return mockResp, nil
		},
	}

	resp, err := mockService.GetTokenFunc(mockReq)

	assert.NoError(t, err)
	assert.Equal(t, mockResp, resp)
}
