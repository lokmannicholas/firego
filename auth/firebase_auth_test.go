package auth

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lokmannicholas/firego"
)

func TestFireAuthImpl_GetFireUser(t *testing.T) {
	if err := firebase.Init(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")); err != nil {
		t.Fatal(err)
	}

	users, err := GetFireAuth().GetAllUsers(5)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", users)
	user, err := GetFireAuth().CreateFirebaseUserByEmail("test@lokmang.com", "1234567")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", user)
	var updateParam UpdateFirebaseUserInfoParam
	err = json.Unmarshal([]byte(`{"username":"user username"}`), &updateParam)
	if err != nil {
		t.Fatal(err)
	}
	updateParam.UID = user.UID
	user, err = GetFireAuth().UpdateFirebaseUser(&updateParam)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", user)

}
