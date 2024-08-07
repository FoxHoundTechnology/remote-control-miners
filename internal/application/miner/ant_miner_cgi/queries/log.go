package queries

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"
)

type LogResponse struct {
	Log string
}

func AntMinerCGILog(clientConnection *http_auth.DigestTransport, username, password, ipAddress string) (*LogResponse, error) {

	newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/log.cgi", ipAddress), nil)
	if err != nil {
		return nil, err
	}

	resp, err := clientConnection.RoundTrip(newRequest)
	if err != nil {
		log.Println("Error in AntMinerCGILog: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in AntMinerCGILog: ", err)
		return nil, err
	}

	return &LogResponse{
		Log: string(body),
	}, nil
}
