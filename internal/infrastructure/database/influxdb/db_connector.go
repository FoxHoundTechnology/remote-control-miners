package database

import (
	"os"

	influxDB "github.com/influxdata/influxdb-client-go/v2"
)

// TODO; batch size based on # of host cores
// TODO: Handle fatal errors

type InfluxDBConnectionSettings struct {
	Client influxDB.Client
	Org    string
	Bucket string
}

// TODO: automate instantiation with init
func Init() *InfluxDBConnectionSettings {

	org := os.Getenv("INFLUXDB_ORG")
	bucket := os.Getenv("INFLUXDB_BUCKET")
	url := "http://influxdb:8086" // NOTE: path has to be identical to container service name
	token := os.Getenv("INFLUXDB_TOKEN")

	client := influxDB.NewClientWithOptions(url, token,
		influxDB.DefaultOptions())

	return &InfluxDBConnectionSettings{
		Client: client,
		Org:    org,
		Bucket: bucket,
	}
}

func Close(i *InfluxDBConnectionSettings) {
	i.Client.Close()
}
