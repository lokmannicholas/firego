package fcm

import (
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
)

type FireMessage struct {
	Fm    FCM
	Title string
	Body  string
}

var channelMessage = make(chan FireMessage, 100)

func Send(m FireMessage) {
	channelMessage <- m
}

func CloudMessagerService() {
	for {
		m := <-channelMessage

		f := &FirebaseCloudMessage{
			Title:   m.Title,
			Body:    m.Body,
			Icon:    "",
			Message: &messaging.Message{},
		}
		switch m.Fm.Device {
		case "IOS":
			err := f.SetAPNMessage(nil, 1)
			if err != nil {
				_ = fmt.Errorf("%+v", err)
			}
		case "ANDROID":
			dur := time.Minute
			err := f.SetAndroidMessage(Normal, "#ffffff", &dur)
			if err != nil {
				_ = fmt.Errorf("%+v", err)
			}

		}
		err := GetFireMessager().SendToDevice(f, m.Fm.Token)
		if err != nil {
			_ = fmt.Errorf("%+v", err)
		}
	}

}
