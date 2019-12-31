package firebase

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"firebase.google.com/go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var ctx = context.Background()
var fireApp *firebase.App
var c *firebase.Config
var credentials *google.Credentials

func Init(path string) error{

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("config file is missing: %v\n", err)
		return  err
	}
	credentials, err = google.CredentialsFromJSON(context.Background(), dat, "https://www.googleapis.com/auth/firebase.remoteconfig")
	if err != nil {
		log.Fatalf("credentials is not found: %v\n", err)
		return  err
	}

	jmap := map[string]string{}
	err = json.Unmarshal(dat, &jmap)
	if err != nil {
		log.Fatalf("config file is not in correct formate: %v\n", err)
		return  err
	}
	config := &firebase.Config{
		ProjectID:     jmap["project_id"],
		StorageBucket: jmap["project_id"] + ".appspot.com"}
	// Initialize default app
	app, err := firebase.NewApp(ctx, config,
		option.WithCredentialsFile(path),
	)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return  err
	}
	c = config
	fireApp = app
	return nil
}
func GetCredentialToken() (string, error) {
	token, err := credentials.TokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}
func GetConfig() *firebase.Config {
	return c
}

func GetFireApp() *firebase.App {
	return fireApp
}
