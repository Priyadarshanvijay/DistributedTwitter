package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/twitter/models"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	// VerifyToken verifies if a given token is correct, if yes, returns the username
	VerifyToken(token string) (*models.User, error)
	// GenerateToken generates a token for a given user
	GenerateToken(userData *models.User) (string, error)
	// ValidateLogin validates the stored password against the password supplied while login attempt
	ValidateLogin(providedData *models.User, trueData *models.User) (*models.User, error)
	// SecureValue returns a secure implementation of the string which can be validated later using the VerifySecureValue function
	SecureValue(value string) string
	// VerifySecureValue returns true if the supplied plain text and secure value are similar
	VerifySecureValue(securedString string, clearTextString string) bool
}

type AuthService struct {
	tokenValidityTime time.Duration
	secretKey         []byte
}

func (a *AuthService) VerifyToken(token string) (*models.User, error) {
	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})
	if err != nil {
		return nil, errors.New("Unauthorized")
	}
	if !parsedToken.Valid {
		return nil, errors.New("Invalid Token")
	}
	return &models.User{
		UserName: claims.Subject,
	}, nil
}

func (a *AuthService) GenerateToken(userData *models.User) (string, error) {
	expiryTime := time.Now().Add(a.tokenValidityTime)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiryTime),
		Subject:   userData.UserName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthService) ValidateLogin(providedData *models.User, trueData *models.User) (*models.User, error) {
	if providedData.UserName != trueData.UserName {
		return nil, errors.New("Invalid login credentials")
	}
	if !a.VerifySecureValue(trueData.UserPassword, providedData.UserPassword) {
		return nil, errors.New("Invalid login credentials")
	}
	return trueData, nil
}

func (a *AuthService) SecureValue(value string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedBytes)
}

func (a *AuthService) VerifySecureValue(hashedString string, clearTextString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(clearTextString))
	return err == nil
}

func New(tokenValidityTime time.Duration, secretKey string) Service {
	authObject := &AuthService{
		tokenValidityTime: tokenValidityTime,
		secretKey:         []byte(secretKey),
	}
	return authObject
}
