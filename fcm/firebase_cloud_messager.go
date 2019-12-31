package fcm

import (
	"context"
	"log"

	"fmt"

	"strings"

	"firebase.google.com/go/messaging"
	"github.com/lokmannicholas/firego"
)

var ctx = context.Background()

type FireMessager interface {
	SendDryRunMode(fcm *FirebaseCloudMessage, token string) error
	SendTopic(fcm *FirebaseCloudMessage, topic ...string) error
	SendToDevice(fcm *FirebaseCloudMessage, fcmToken string) error
}

type FireMessagerImpl struct {
	Messager *messaging.Client
}

func (b *FireMessagerImpl) SendDryRunMode(fcm *FirebaseCloudMessage, token string) error {
	fcm.Token = token
	response, err := b.Messager.SendDryRun(ctx, fcm.Build())
	if err != nil {
		log.Fatalln(err)
		return err
	}
	// Response is a message ID string.
	fmt.Println("Dry run successful:", response)
	return nil
}
func (b *FireMessagerImpl) SendTopic(fcm *FirebaseCloudMessage, topic ...string) error {
	if len(topic) == 1 && len(topic) > 0 {
		fcm.Topic = topic[0]
	} else {
		conds := []string{}
		for _, s := range topic {
			conds = append(conds, fmt.Sprintf(`'%s' in topics `, s))
		}
		condition := strings.Join(conds, " || ")
		fcm.Condition = condition
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := b.Messager.Send(ctx, fcm.Build())
	if err != nil {
		log.Fatalln(err)
		return err
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	return nil
}
func (b *FireMessagerImpl) SendToDevice(fcm *FirebaseCloudMessage, fcmToken string) error {

	fcm.Token = fcmToken
	// registration token.
	response, err := b.Messager.Send(ctx, fcm.Build())
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetFireMessager() FireMessager {

	//storage bucket
	client, err := firebase.GetFireApp().Messaging(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	return &FireMessagerImpl{
		Messager: client,
	}
}
