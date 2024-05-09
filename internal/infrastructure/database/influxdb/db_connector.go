package database

import (
	"os"

	influxDB "github.com/influxdata/influxdb-client-go/v2"
)

// TODO; batch size based on # of host cores
// TODO: Handle fatal errors
// TODO: env

type InfluxDBConnectionSettings struct {
	Client influxDB.Client
	Org    string
	Bucket string
}

// TODO: automate instantiation with init
func Init() InfluxDBConnectionSettings {

	org := os.Getenv("INFLUX_DB_ORG_NAME")
	bucket := os.Getenv("INFLUX_DB_BUCKET_NAME")
	url := "http://influxdb:8086" // = container name
	token := os.Getenv("INFLUX_DB_TOKEN")

	client := influxDB.NewClientWithOptions(url, token,
		influxDB.DefaultOptions().SetBatchSize(1000))

	return InfluxDBConnectionSettings{
		Client: client,
		Org:    org,
		Bucket: bucket,
	}
}

func Close(i *InfluxDBConnectionSettings) {
	i.Client.Close()
}
