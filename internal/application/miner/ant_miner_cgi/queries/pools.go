package queries

import (
	"encoding/json"
	"fmt"
	"foxhound/pkg/http_auth"
	"io"
	"net/http"
)

// cgi-bin/pools.cgi: Get pool information.
// TODO: default values for response object

type rawGetPoolsResponse struct {
	Status Status `json:"STATUS"`
	Pools  []Pool `json:"POOLS"`
}

type Pool struct {
	Index     int    `json:"index"`  // pool index
	URL       string `json:"url"`    // pool url
	User      string `json:"user"`   // miner name = pool user name
	Status    string `json:"status"` // pool status : "Alive"
	Priority  int    `json:"priority"`
	Getworks  int    `json:"getworks"`
	Accepted  int    `json:"accepted"` // accepted shares
	Rejected  int    `json:"rejected"` // rejected shares
	Discarded int    `json:"discarded"`
	Stale     int    `json:"stale"` // stale shares
	Diff      string `json:"diff"`
	Diff1     int    `json:"diff1"`
	Diffa     int    `json:"diffa"`
	Diffr     int    `json:"diffr"`
	Diffs     int    `json:"diffs"`
	Lsdiff    int    `json:"lsdiff"`
	Lstime    string `json:"lstime"`
}

type GetPoolsResponse struct {
	Index    int    `json:"index"`     // pool index
	URL      string `json:"url"`       // pool url
	UserName string `json:"user_name"` // miner name = pool user name
	Password string `json:"password"`  // pool password
	Status   string `json:"status"`    // pool status : "Alive"
	Accepted int    `json:"accepted"`  // accepted shares
	Rejected int    `json:"rejected"`  // rejected shares
	Stale    int    `json:"stale"`     // stale shares
}

func AntMinerCGIGetPools(username, password, ipAddress string) (*[]GetPoolsResponse, error) {

	t := http_auth.NewTransport(username, password)

	newRequest, err := http.NewRequest("GET", fmt.Sprintf("http://%s/cgi-bin/pools.cgi", ipAddress), nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.RoundTrip(newRequest)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawGetPoolsResponse rawGetPoolsResponse
	err = json.Unmarshal(body, &rawGetPoolsResponse)
	if err != nil {
		return nil, err
	}

	fmt.Println("rawGetPoolsResponse: ", rawGetPoolsResponse.Pools[0].Index, rawGetPoolsResponse.Pools[0].URL, rawGetPoolsResponse.Pools[0].User, rawGetPoolsResponse.Pools[0].Status, rawGetPoolsResponse.Pools[0].Accepted, rawGetPoolsResponse.Pools[0].Rejected, rawGetPoolsResponse.Pools[0].Stale)

	var pools []GetPoolsResponse
	for _, pool := range rawGetPoolsResponse.Pools {
		pools = append(pools, GetPoolsResponse{
			Index:    pool.Index,
			URL:      pool.URL,
			UserName: pool.User,
			Status:   pool.Status,
			Accepted: pool.Accepted,
			Rejected: pool.Rejected,
			Stale:    pool.Stale,
		})
	}

	return &pools, nil
}
