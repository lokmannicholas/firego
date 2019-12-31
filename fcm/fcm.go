package fcm

import (
	"errors"
	"time"

	"encoding/json"

	"firebase.google.com/go/messaging"
)

type FirebaseCloudMessage struct {
	*messaging.Message
	Title string
	Body  string
	Icon  string
}

type AndroidPriority string

const Normal AndroidPriority = "normal"

func (f *FirebaseCloudMessage) SetAndroidMessage(priority AndroidPriority, color string, duration *time.Duration) error {
	f.Android = &messaging.AndroidConfig{
		TTL:      duration,
		Priority: string(priority),
		Notification: &messaging.AndroidNotification{
			Title: f.Title,
			Body:  f.Body,
			Icon:  f.Icon,
			Color: color,
		},
	}
	return nil
}

func (f *FirebaseCloudMessage) SetWebMessage() error {
	f.Webpush = &messaging.WebpushConfig{
		Notification: &messaging.WebpushNotification{
			Title: f.Title,
			Body:  f.Body,
			Icon:  f.Icon,
		},
	}
	return nil
}

//Ref: https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html
type APNHeader struct {
	ID         string `json:"apns-id"`
	CollapseID string `json:"apns-collapse-id"`
	Expiration string `json:"apns-expiration"`
	Topic      string `json:"apns-topic"`
	Priority   string `json:"apns-priority"`
}

func (f *FirebaseCloudMessage) SetAPNMessage(apnHeader *APNHeader, badge int) error {

	header := map[string]string{}
	if apnHeader != nil {
		data, err := json.Marshal(apnHeader)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &header)
		if err != nil {
			return err
		}
	}

	f.APNS = &messaging.APNSConfig{
		Headers: header,
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: f.Title,
					Body:  f.Body,
				},
				Badge: &badge,
			},
		},
	}
	return nil
}
func (f *FirebaseCloudMessage) Build() *messaging.Message {
	m := &messaging.Message{
		Android: f.Android,
		Webpush: f.Webpush,
		APNS:    f.APNS,
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
	}
	if len(f.Token) != 0 {
		m.Token = f.Token
		return m
	}
	if len(f.Condition) != 0 {
		m.Condition = f.Condition
		return m
	}
	if len(f.Topic) != 0 {
		m.Topic = f.Topic
		return m
	}
	return m
}

type FCM struct {
	Device string
	Token  string
}

func ToFCM(mod interface{}) (*FCM, error) {
	if mod == nil {
		return nil, errors.New("input data invalid")
	}
	return mod.(*FCM), nil
}
