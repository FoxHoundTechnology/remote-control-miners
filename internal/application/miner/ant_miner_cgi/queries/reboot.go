package queries

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/FoxHoundTechnology/remote-control-miners/foxhound/pkg/http_auth"

	"github.com/sirupsen/logrus"
)

// cgi-bin/reboot.cgi: Reboot the miner.

func AntMinerCGIReboot(username, password, ipAddress string) error {

	t := http_auth.NewTransport(username, password)

	newRequest, err := http.NewRequest("GET", fmt.Sprintf("http://%s/cgi-bin/reboot.cgi", ipAddress), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"newRequest": newRequest,
		}).Debug("Error creating new request")
		return err
	}

	resp, err := t.RoundTrip(newRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"resp":  resp,
		}).Debug("Error creating new request")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"body":  body,
		}).Debug("Error reading response body")
		return err
	}

	log.Println("RESULT OF REBOOT", string(body))

	return nil
}
