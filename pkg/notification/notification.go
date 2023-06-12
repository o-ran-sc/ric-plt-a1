package notification

import (
	"bytes"
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
)

func SendNotification(notificationDestination string, body string) error {
	a1.Logger.Debug("In SendNotification")
	client := &http.Client{}
	req, err := http.NewRequest("POST", notificationDestination, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "a1mediator")
	res, err := client.Do(req)
	if err != nil {
		a1.Logger.Debug("Error sending POST message to %v+", notificationDestination)
		return err
	}
	a1.Logger.Debug("Notification response %s received from far-end", res.Status)
	return nil
}
