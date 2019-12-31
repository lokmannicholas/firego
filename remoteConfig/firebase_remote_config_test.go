package remoteConfig

import (
	"os"
	"testing"

	"github.com/lokmannicholas/firego"
)

func TestFirebaseRemoteConfig_Get(t *testing.T) {
	if err := firebase.Init(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")); err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	err := GetRemoteConfig(&v)
	if err != nil {
		t.Error(err)
	}
	t.Log(v)
}
