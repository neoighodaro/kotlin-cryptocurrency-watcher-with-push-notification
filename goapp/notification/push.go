package notification

import (
	"fmt"

	"github.com/pusher/push-notifications-go"
)

const (
	instanceID = "YOUR_INSTANCE_ID_HERE"
	secretKey  = "YOUR_SECRET_KEY_HERE"
)

// SendNotification sends push notification to devices
func SendNotification(currency string, price float64, uuid string) error {
	notifications, err := pushnotifications.New(instanceID, secretKey)
	if err != nil {
		return err
	}

	publishRequest := map[string]interface{}{
		"fcm": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": currency + " Price Change",
				"body":  fmt.Sprintf("The price of %s has changed to %f", currency, price),
			},
		},
	}

	interest := fmt.Sprintf("%s_%s_changed", uuid, currency)
	pubID, err := notifications.Publish([]string{interest}, publishRequest)
	if err != nil {
		return err
	}

	println(pubID)

	return nil
}
