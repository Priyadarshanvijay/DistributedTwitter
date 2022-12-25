package auth

import (
	"testing"
	"time"

	"github.com/twitter/models"
)

const (
	Test_DefaultTokenValidityTime = 4 * time.Second
	Test_DefaultSecretString      = "secret?"
)

func TestAuthService_VerifySecureValue(t *testing.T) {
	test_auth_service := New(Test_DefaultTokenValidityTime, Test_DefaultSecretString)
	og_value := "this_is_test_value"
	not_og_value := "this_is_test_value1"
	secured_val := test_auth_service.SecureValue(og_value)
	ok := test_auth_service.VerifySecureValue(secured_val, og_value)
	if !ok {
		t.Error("Unable to verify secured value")
	}
	ok = test_auth_service.VerifySecureValue(secured_val, not_og_value)
	if ok {
		t.Error("Error: matching with incorrect clear text value")
	}
	rev_ok := test_auth_service.VerifySecureValue(og_value, secured_val)
	if rev_ok {
		t.Error("Error: reversed args giving ok")
	}
}

func TestAuthService_VerifyToken(t *testing.T) {
	test_auth_service := New(Test_DefaultTokenValidityTime, Test_DefaultSecretString)
	non_token := ""
	_, err := test_auth_service.VerifyToken(non_token)
	if err == nil {
		t.Error("Invalid token not returning error")
	} else if err.Error() != "Unauthorized" {
		t.Errorf("Empty token giving non-Empty error: %s\n", err.Error())
	}
	test_user := &models.User{
		UserName: "test_user",
	}
	invalid_token, err := New(Test_DefaultTokenValidityTime, "not_secret?").GenerateToken(test_user)
	_, err = test_auth_service.VerifyToken(invalid_token)
	if err == nil {
		t.Error("Invalid token not returning error")
	} else if err.Error() != "Unauthorized" {
		t.Errorf("Unexpected error: %s\n", err.Error())
	}
	expired_token, err := New(1, Test_DefaultSecretString).GenerateToken(test_user)
	time.Sleep(2)
	_, err = test_auth_service.VerifyToken(expired_token)
	if err == nil {
		t.Error("Invalid token not returning error")
	} else if err.Error() != "Unauthorized" {
		t.Errorf("Unexpected error: %s\n", err.Error())
	}
	valid_token, err := test_auth_service.GenerateToken(test_user)
	decoded_user, err := test_auth_service.VerifyToken(valid_token)
	if err != nil {
		t.Errorf("Valid token returning error: %+v\n", err)
	}
	if (decoded_user != nil) && (decoded_user.UserName != test_user.UserName) {
		t.Errorf("Valid token returning incorrect username: %+v\n", decoded_user)
	}
}

func TestAuthService_ValidateLogin(t *testing.T) {
	test_auth_service := New(Test_DefaultTokenValidityTime, Test_DefaultSecretString)
	test_user := &models.User{
		UserName:     "test_user",
		UserPassword: "password@1",
	}
	stored_user := &models.User{
		UserName:     "test_user",
		UserPassword: test_auth_service.SecureValue(test_user.UserPassword),
	}
	returned_user, err := test_auth_service.ValidateLogin(test_user, stored_user)
	if err != nil {
		t.Errorf("Unexpected error in correct login: %+v\n", err)
	}
	if returned_user == nil {
		t.Errorf("Nil user returned")
	}
	if returned_user != nil && returned_user.UserPassword != stored_user.UserPassword {
		t.Errorf("Incorrect user returned: %+v\n", returned_user)
	}
	wrong_username := &models.User{
		UserName:     "not_test_user",
		UserPassword: "password@1",
	}
	wrong_password := &models.User{
		UserName:     "test_user",
		UserPassword: "password@2",
	}
	returned_user, err = test_auth_service.ValidateLogin(wrong_username, stored_user)
	if err == nil {
		t.Errorf("Nil error returned for wrong username\n")
	}
	if returned_user != nil {
		t.Errorf("Incorrect user returned in wrong username: %+v\n", returned_user)
	}
	if err != nil && err.Error() != "Invalid login credentials" {
		t.Errorf("Incorrect error returned in wrong username: %+v\n", err)
	}
	returned_user, err = test_auth_service.ValidateLogin(wrong_password, stored_user)
	if err == nil {
		t.Errorf("Nil error returned for wrong password\n")
	}
	if returned_user != nil {
		t.Errorf("Incorrect user returned in wrong password: %+v\n", returned_user)
	}
	if err != nil && err.Error() != "Invalid login credentials" {
		t.Errorf("Incorrect error returned in wrong password: %+v\n", err)
	}
}
