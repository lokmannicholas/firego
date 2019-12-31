package fcm

import (
	"os"
	"testing"

	"github.com/lokmannicholas/firego"
)

func TestFCM(t *testing.T) {
	if err := firebase.Init(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")); err != nil {
		t.Fatal(err)
	}
	go CloudMessagerService()
	f := FireMessage{
		Title: "FCM Test",
		Body:  "This is test",
	}
	f.Fm = FCM{
		Device: "IOS",
		Token:  "12345678",
	}
	go Send(f)
}
