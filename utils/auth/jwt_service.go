package auth

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(UserID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey []byte
}

func NewService() *jwtService {
	jwtSecret := os.Getenv("JWT_SECRET")

	return &jwtService{
		secretKey: []byte(jwtSecret),
	}
}

func (s *jwtService) GenerateToken(UserID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token")
		}

		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
