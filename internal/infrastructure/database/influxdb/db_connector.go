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
func Init() InfluxDBConnectionSettings {

	org := os.Getenv("INFLUX_DB_ORG")
	bucket := os.Getenv("INFLUX_DB_BUCKET")
	url := os.Getenv("INFLUX_DB_URL")
	port := os.Getenv("INFLUX_DB_PORT")
	path := url + ":" + port
	token := os.Getenv("INFLUX_DB_TOKEN")

	client := influxDB.NewClientWithOptions(path, token,
		influxDB.DefaultOptions().SetBatchSize(10000))

	return InfluxDBConnectionSettings{
		Client: client,
		Org:    org,
		Bucket: bucket,
	}
}

func Close(i *InfluxDBConnectionSettings) {
	i.Client.Close()
}
