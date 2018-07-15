package notification

import (
	"fmt"
	"strconv"

	"github.com/pusher/push-notifications-go"
)

const (
	instanceID = "PUSHER_BEAMS_INSTANCE_ID"
	secretKey  = "PUSHER_BEAMS_SECRET_KEY"
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
				"body":  fmt.Sprintf("The price of %s has changed to $%s", currency, strconv.FormatFloat(price, 'f', 2, 64)),
			},
		},
	}

	interest := fmt.Sprintf("%s_%s_changed", uuid, currency)
	_, err = notifications.Publish([]string{interest}, publishRequest)
	if err != nil {
		return err
	}

	return nil
}
