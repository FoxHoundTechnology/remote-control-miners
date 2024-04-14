package queries

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
	Status   string `json:"status"`    // pool status : "Alive"
	Accepted int    `json:"accepted"`  // accepted shares
	Rejected int    `json:"rejected"`  // rejected shares
	Stale    int    `json:"stale"`     // stale shares
}

func AntMinerCGIGetPools(ipAddress string) (*GetPoolsResponse, error) {
	return nil, nil
}
