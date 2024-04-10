package scanner

type MinerType string

const (
	AntminerCgi MinerType = "antminer_cgi"
	//...
)

type Status struct {
	Active bool
}

type Config struct {
	Interval int // in minutes
	Username string
	Password string
}
