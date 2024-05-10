package repositories

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	influxDB "github.com/influxdata/influxdb-client-go/v2"
	influxDB_api "github.com/influxdata/influxdb-client-go/v2/api"

	timeseries_database "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/influxdb"
)

// TODO: context cancellation
// TODO: data race condition
// TODO: RW mutex

type MinerTimeSeriesRepository struct {
	db     timeseries_database.InfluxDBConnectionSettings
	writer influxDB_api.WriteAPI

	rw *sync.RWMutex

	timeseriesMinerData []MinerTimeSeries
	timeseriesPoolData  []PoolTimeSeries
}

func NewMinerTimeSeriesRepository(db timeseries_database.InfluxDBConnectionSettings) *MinerTimeSeriesRepository {

	return &MinerTimeSeriesRepository{
		db: db,

		writer: db.Client.WriteAPI(db.Org, db.Bucket),
		rw:     new(sync.RWMutex),

		timeseriesMinerData: []MinerTimeSeries{},
		timeseriesPoolData:  []PoolTimeSeries{},
	}
}

func (r *MinerTimeSeriesRepository) WriteMinerData(mac_address string, data MinerTimeSeries) error {

	r.rw.Lock()
	defer r.rw.Unlock()

	r.timeseriesMinerData = append(r.timeseriesMinerData, data)

	return nil
}

func (r *MinerTimeSeriesRepository) FlushMinerData() error {

	fmt.Println("flushing miner data with length", len(r.timeseriesMinerData))

	for _, data := range r.timeseriesMinerData {

		temperatureStringArray := make([]string, len(data.TempSensor))
		for index, temperature := range data.TempSensor {
			temperatureStringArray[index] = fmt.Sprintf("%d", temperature)
		}

		fanStringArray := make([]string, len(data.FanSensor))
		for index, fan_speed := range data.FanSensor {
			fanStringArray[index] = fmt.Sprintf("%d", fan_speed)
		}

		fields := map[string]interface{}{
			"hashrate":     data.HashRate,
			"temp_sensors": strings.Join(temperatureStringArray, ","),
			"fan_sensors":  strings.Join(fanStringArray, ","),
		}

		tag := map[string]string{
			"mac_address": data.MacAddress,
		}

		point := influxDB.NewPoint(
			"miner_data",
			tag,
			fields,
			time.Now(),
		)
		r.writer.WritePoint(point)
	}

	r.writer.Flush()
	r.timeseriesMinerData = nil

	return nil
}

func (r *MinerTimeSeriesRepository) WritePoolData(mac_address string, data PoolTimeSeries) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.timeseriesPoolData = append(r.timeseriesPoolData, data)
	return nil
}

func (r *MinerTimeSeriesRepository) FlushPoolData() error {

	for _, data := range r.timeseriesPoolData {

		fields := map[string]interface{}{
			"accepted": data.Accepted,
			"rejected": data.Rejected,
			"stale":    data.Stale,
		}

		tag := map[string]string{
			"mac_address": data.MacAddress,
		}

		point := influxDB.NewPoint(
			"pool_data",
			tag,
			fields,
			time.Now(),
		)
		r.writer.WritePoint(point)
	}

	r.writer.Flush()
	r.timeseriesPoolData = nil
	return nil
}

// NOTE: mac_address is null in the response object
func (r *MinerTimeSeriesRepository) ReadMinerData(mac_address string, interval int) (MinerTimeSeriesResponse, error) {
	queryAPI := r.db.Client.QueryAPI(r.db.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -%dm)
	|> filter(fn: (r) => r._measurement == "miner_data" and r.mac_address == "%s")
	|> sort(columns: ["_time"], desc: false)`,
		r.db.Bucket, interval, mac_address)

	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return MinerTimeSeriesResponse{}, err
	}

	minerDataMapArray := make(map[time.Time]*MinerTimeSeries)
	var timeStamps []time.Time

	for results.Next() {

		t := results.Record().Time()

		if _, exists := minerDataMapArray[t]; !exists {
			timeStamps = append(timeStamps, t)
			minerDataMapArray[t] = &MinerTimeSeries{}
		}

		minerData := minerDataMapArray[t]
		var hashrate int
		var sensorData string // NOTE: temp_sensors or fan_sensors
		fieldName := results.Record().Field()

		switch v := results.Record().Value().(type) {
		case int64:
			hashrate = int(v)
		case float64:
			hashrate = int(v)
		case string:
			sensorData = v
		default:
			fmt.Println("unknown type")
		}

		switch fieldName {
		case "hashrate":
			minerData.HashRate = hashrate
		case "temp_sensors":

			// perhaps updateing the array here?

			temperatureStringArray := strings.Split(sensorData, ",")
			temperatureSlice := make([]int, len(temperatureStringArray))

			for index, temperatureString := range temperatureStringArray {
				temperatureValue, err := strconv.Atoi(temperatureString)
				if err != nil {
					fmt.Printf("error converting temperature value: %s\n", err)
				}
				temperatureSlice[index] = temperatureValue
			}

			minerData.TempSensor = temperatureSlice

		case "fan_sensors":
			fanStringArray := strings.Split(sensorData, ",")
			fanSlice := make([]int, len(fanStringArray))

			fmt.Println("FAN STRING ARRAY", fanStringArray)

			for index, fanString := range fanStringArray {
				fanValue, err := strconv.Atoi(fanString)
				if err != nil {
					fmt.Printf("error converting fan value: %s\n", err)
				}
				fanSlice[index] = fanValue
			}

			minerData.FanSensor = fanSlice

		}
	}

	if err := results.Err(); err != nil {
		return MinerTimeSeriesResponse{}, fmt.Errorf("error in response: %v", err)
	}

	// sorting the timestamp here
	sort.Slice(timeStamps, func(i, j int) bool {
		return timeStamps[i].Before(timeStamps[j])
	})

	var minerTimeSeriesArray []MinerTimeSeries

	// reordering the miner data based on the sorted timestamp order
	for _, timestamp := range timeStamps {
		minerTimeSeriesArray = append(minerTimeSeriesArray, *minerDataMapArray[timestamp])
	}

	return MinerTimeSeriesResponse{
		Record:     minerTimeSeriesArray,
		TimeStamps: timeStamps,
	}, nil
}

func (r *MinerTimeSeriesRepository) ReadPoolData(mac_address string, interval int) (PoolTimeSeriesResponse, error) {
	queryAPI := r.db.Client.QueryAPI(r.db.Org)

	// Modify the range to use the interval for days.
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -%dm) 
		|> filter(fn: (r) => r._measurement == "pool_data" and r.mac_address == "%s")
		|> sort(columns: ["_time"], desc: false)`,
		r.db.Bucket, interval, mac_address)

	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return PoolTimeSeriesResponse{}, fmt.Errorf("error constructing the request query: %v", err)
	}

	poolDataMapArray := make(map[time.Time]*PoolTimeSeries)
	var timeStamps []time.Time

	for results.Next() {
		t := results.Record().Time()

		if _, exists := poolDataMapArray[t]; !exists {
			timeStamps = append(timeStamps, t)
			poolDataMapArray[t] = &PoolTimeSeries{}
		}

		var value int
		fieldName := results.Record().Field()

		switch v := results.Record().Value().(type) {
		case int64:
			value = int(v)
		case float64:
			value = int(v)
		default:
			fmt.Println("unknown type")
		}

		fmt.Println("THIS IS THE VALUE INSERING INTO THE POOL DATA MAP ARRAY", value)

		poolData := poolDataMapArray[t]
		switch fieldName {
		case "accepted":
			poolData.Accepted = value
		case "rejected":
			poolData.Rejected = value
		case "stale":
			poolData.Stale = value

		}
	}

	if err := results.Err(); err != nil {
		return PoolTimeSeriesResponse{}, fmt.Errorf("error in response: %v", err)
	}

	// sorting the timestampArray here
	sort.Slice(timeStamps, func(i, j int) bool {
		return timeStamps[i].Before(timeStamps[j])
	})

	var poolTimeSeriesArray []PoolTimeSeries

	for _, timestamp := range timeStamps {
		poolTimeSeriesArray = append(poolTimeSeriesArray, *poolDataMapArray[timestamp])
	}

	return PoolTimeSeriesResponse{
		Record:     poolTimeSeriesArray,
		TimeStamps: timeStamps,
	}, nil
}
