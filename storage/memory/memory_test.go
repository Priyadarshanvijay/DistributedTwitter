package memory

import (
	"reflect"
	"testing"

	"github.com/twitter/models"
)

func Test_userStore_AddUser(t *testing.T) {
	test_user := &models.User{
		UserName:     "test1",
		UserPassword: "password1",
	}
	test_storage := New()
	created_user, err := test_storage.UserStore().AddUser(test_user)
	if err != nil {
		t.Errorf("Error in adding user: %+v\n", err)
	}
	if created_user.UserName != test_user.UserName {
		t.Error("Returned Users username doesn't match")
	}
	if created_user.UserPassword != test_user.UserPassword {
		t.Error("Returned Users username doesn't match")
	}

	created_user_v2, err := test_storage.UserStore().AddUser(test_user)

	if err == nil {
		t.Error("Error in adding user, able to add duplicate users")
	}
	if err.Error() != "User already exists" {
		t.Errorf("Error in adding duplicate user, unexpected message: %+v\n", err)
	}
	if created_user_v2 != nil {
		t.Errorf("Error in adding duplicate user, returned non nil user: %+v\n", created_user_v2)
	}
}

func Test_userStore_GetUser(t *testing.T) {
	test_storage := New()
	test_user := &models.User{
		UserName:     "test1",
		UserPassword: "password1",
	}
	nil_user, err := test_storage.UserStore().GetUser(test_user.UserName)

	if err == nil {
		t.Error("Error in get nil user, user was never created")
	}

	if err.Error() != "User doesn't exists" {
		t.Errorf("Error in get nil user, unexpected message: %+v\n", err)
	}

	if nil_user != nil {
		t.Errorf("Error in getting nil user, returned non nil user: %+v\n", nil_user)
	}

	test_storage.UserStore().AddUser(test_user)

	retrieved_user, err := test_storage.UserStore().GetUser(test_user.UserName)

	if err != nil {
		t.Errorf("Error in get user: %+v\n", err)
	}

	if retrieved_user.UserName != test_user.UserName {
		t.Error("Error in get user: Username doesn't match")
	}

	if retrieved_user.UserPassword != test_user.UserPassword {
		t.Error("Error in get user: UserPassword doesn't match")
	}

}

func Test_userStore_UpdateUser(t *testing.T) {
	test_storage := New()
	username := "test1"
	test_user_init := &models.User{
		UserName:     username,
		UserPassword: "password1",
		UserEmail:    "old_email@email.com",
	}
	test_user_updated_password := &models.User{
		UserName:     username,
		UserPassword: "newPassword",
	}

	test_user_updated_email := &models.User{
		UserName:  username,
		UserEmail: "new_email@email.com",
	}

	test_storage.UserStore().AddUser(test_user_init)

	retrieved_user, _ := test_storage.UserStore().GetUser(username)

	if !reflect.DeepEqual(test_user_init, retrieved_user) {
		t.Error("User did not match before updating")
	}

	test_storage.UserStore().UpdateUser(test_user_updated_password)

	updated_password_user, _ := test_storage.UserStore().GetUser(username)

	if updated_password_user.UserEmail != test_user_init.UserEmail {
		t.Error("UserEmail changed unexpectedly on password update")
	}

	if updated_password_user.UserPassword != test_user_updated_password.UserPassword {
		t.Error("UserPassword did not change correctly on password update")
	}

	test_storage.UserStore().UpdateUser(test_user_updated_email)

	updated_email_user, _ := test_storage.UserStore().GetUser(username)

	if updated_email_user.UserPassword != test_user_updated_password.UserPassword {
		t.Error("UserPassword changed unexpectedly on email update")
	}

	if updated_email_user.UserEmail != test_user_updated_email.UserEmail {
		t.Error("UserEmail did not change correctly on email update")
	}
}
