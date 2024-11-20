package firebase_utils

import (
	"firebase.google.com/go/v4/messaging"
)

func BuildTopicMessage(
	title, body, topic *string,
	data map[string]string,
) *messaging.Message {
	message := &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title: *title,
			Body:  *body,
			// ImageURL: "",
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Sound: "default",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
				},
			},
		},
		Topic: *topic,
	}

	return message
}

func BuildTokenMessage(
	title, body, userFCMToken *string,
	data map[string]string,
) *messaging.Message {
	message := &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title: *title,
			Body:  *body,
			// ImageURL: "",
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Sound: "default",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
				},
			},
		},
		Token: *userFCMToken,
	}
	return message
}
