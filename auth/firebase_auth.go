package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/auth"
	"github.com/lokmannicholas/firego"
	"google.golang.org/api/iterator"
)

type FireAuth interface {
	GetFirebaseUser(schema GetFirebaseUserSchema) *FirebaseUserInfo
	GetAllUsers(userNumber int) ([]FirebaseUserInfo, error)
	CheckFireToken(UID string, idToken string) (*auth.Token, error)
	CustomToken(UID string) string
	CreateFirebaseUserByEmail(email, password string) (*FirebaseUserInfo, error)
	UpdateFirebaseUser(fireInfo *UpdateFirebaseUserInfoParam) (*FirebaseUserInfo, error)
	DeleteFirebaseUser(uid string) error
	DeactiveFirebaseUser(UID string) (*FirebaseUserInfo, error)
}

type FirebaseUserInfo struct {
	UID                string
	DisplayName        string `json:"username"`
	Email              string
	PhoneNumber        string `json:"tele"`
	PhotoURL           string
	Password           string
	ProviderID         string
	CreationTimestamp  *time.Time
	LastLogInTimestamp *time.Time
	Provider           string
	Disabled           bool
	EmailVerified      bool
}
type FireAuthImpl struct {
	FirebaseAuth *auth.Client
}

func GetFireAuth() FireAuth {
	//storage bucket
	client, err := firebase.GetFireApp().Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return &FireAuthImpl{
		FirebaseAuth: client,
	}
}
func (b *FireAuthImpl) CustomToken(UID string) string {
	token, err := b.FirebaseAuth.CustomToken(context.Background(), UID)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)
	return token
}

type GetFirebaseUserSchema struct {
	UID   *string
	Email *string
	Phone *string
	Token *string
}

func (b *FireAuthImpl) GetFirebaseUser(schema GetFirebaseUserSchema) *FirebaseUserInfo {
	var err error
	var u *auth.UserRecord
	if schema.UID != nil {
		u, err = b.FirebaseAuth.GetUser(context.Background(), *schema.UID)
	} else if schema.Email != nil {
		u, err = b.FirebaseAuth.GetUserByEmail(context.Background(), *schema.Email)
	} else if schema.Phone != nil {
		u, err = b.FirebaseAuth.GetUserByPhoneNumber(context.Background(), *schema.Phone)
	}else if schema.Token != nil {
		var token *auth.Token
		token, err = b.FirebaseAuth.VerifyIDToken(context.Background(), *schema.Token)
		u, err = b.FirebaseAuth.GetUser(context.Background(), token.UID)
	}

	if err != nil {
		log.Fatalf("error getting user with schema %+v: %v\n", schema, err)
	}
	return ToFirebaseUserInfo(u)
}

func (b *FireAuthImpl) CreateFirebaseUserByEmail(email, password string) (*FirebaseUserInfo, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		//PhoneNumber("+15555550100").
		Password(password).
		DisplayName(email).
		//PhotoURL("http://www.example.com/12345678/photo.png").
		Disabled(false)
	u, err := b.FirebaseAuth.CreateUser(context.Background(), params)
	if err != nil {
		return nil, err
	}
	fu := ToFirebaseUserInfo(u)
	createTime := time.Unix(u.UserMetadata.CreationTimestamp, 0)
	fu.CreationTimestamp = &createTime
	return fu, nil
}

type UpdateFirebaseUserInfoParam struct {
	UID         string
	DisplayName *string `json:"username"`
	Email       *string
	PhoneNumber *string `json:"tele"`
	PhotoURL    *string
	Password    *string
}

func (b *FireAuthImpl) UpdateFirebaseUser(fireInfo *UpdateFirebaseUserInfoParam) (*FirebaseUserInfo, error) {
	params := &auth.UserToUpdate{}
	if fireInfo.Email != nil {
		params.Email(*fireInfo.Email)
	}
	if fireInfo.DisplayName != nil {
		params.DisplayName(*fireInfo.DisplayName)
	}
	if fireInfo.PhoneNumber != nil {
		params.PhoneNumber(*fireInfo.PhoneNumber)
	}
	if fireInfo.Password != nil {
		params.Password(*fireInfo.Password)
	}
	if fireInfo.PhotoURL != nil {
		params.PhotoURL(*fireInfo.PhotoURL)
	}

	u, err := b.FirebaseAuth.UpdateUser(context.Background(), fireInfo.UID, params)
	if err != nil {
		return nil, err
	}
	return ToFirebaseUserInfo(u), nil
}

func (b *FireAuthImpl) DeleteFirebaseUser(uid string) error {
	return b.FirebaseAuth.DeleteUser(context.Background(), uid)
}

func (b *FireAuthImpl) DeactiveFirebaseUser(UID string) (*FirebaseUserInfo, error) {
	params := (&auth.UserToUpdate{}).Disabled(true)

	u, err := b.FirebaseAuth.UpdateUser(context.Background(), UID, params)
	if err != nil {
		return nil, err
	}

	return ToFirebaseUserInfo(u), nil
}

func (b *FireAuthImpl) GetAllUsers(userNumber int) ([]FirebaseUserInfo, error) {
	var users []FirebaseUserInfo
	pager := iterator.NewPager(b.FirebaseAuth.Users(context.Background(), ""), userNumber, "")
	for {

		var authUsers []*auth.ExportedUserRecord
		nextPageToken, err := pager.NextPage(&authUsers)
		if err != nil {
			log.Fatalf("paging error %v\n", err)
			return nil, err
		}
		for _, u := range authUsers {
			users = append(users, *ToFirebaseUserInfo(u.UserRecord))
		}

		if nextPageToken == "" {
			break
		}
	}
	return users, nil
}

func (b *FireAuthImpl) CheckFireToken(UID string, idToken string) (*auth.Token, error) {
	token, err := b.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		_ = fmt.Errorf("error verifying ID token: %v\n", err)
		return nil, err
	}

	if token.UID != UID {
		err := errors.New("user not match with correct token")
		_ = fmt.Errorf("%+v", err)
		return token, err
	}
	return token, nil
}

func ToFirebaseUserInfo(u *auth.UserRecord) *FirebaseUserInfo {
	createTime := time.Unix(u.UserMetadata.CreationTimestamp, 0)
	lastLogin := time.Unix(u.UserMetadata.LastLogInTimestamp, 0)

	return &FirebaseUserInfo{
		UID:                u.UID,
		DisplayName:        u.UserInfo.DisplayName,
		Email:              u.UserInfo.Email,
		PhoneNumber:        u.UserInfo.PhoneNumber,
		PhotoURL:           u.UserInfo.PhotoURL,
		ProviderID:         u.UserInfo.ProviderID,
		CreationTimestamp:  &createTime,
		LastLogInTimestamp: &lastLogin,
		Provider:           u.UserInfo.ProviderID,
		Disabled:           u.Disabled,
		EmailVerified:      u.EmailVerified,
	}
}
