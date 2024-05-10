package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	postgres "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/postgres"

	http_auth "github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	// TODO: db migration/seed
	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	scanner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/scanner/domain"

	fleet_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/fleet"
	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"
	scanner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/scanner"

	timeseries_database "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/influxdb"

	routes "github.com/FoxHoundTechnology/remote-control-miners/internal/interface/routers"

	ant_miner_cgi_queries "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/queries"
	ant_miner_cgi_service "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/service"
)

// TODO: retrieve a collection of max temperature and max fan speed
// TODO: select statement for different vendors
// TODO: batch operation for miner stats update
// TODO: R&D for pool library's memory leak
// TODO: logic for detecting offline miners
// TODO: logic for identifying the active pool
// TODO: logic for combined miner error supports

func main() {

	postgresDB := postgres.Init()
	InfluxDBConnectionSettings := timeseries_database.Init()

	err := postgresDB.AutoMigrate(
		// NOTE: The order matters
		&fleet_repo.Fleet{},

		&scanner_repo.Scanner{},
		&scanner_repo.Alert{},
		&scanner_repo.AlertCondition{},
		&scanner_repo.AlertLog{},

		&miner_repo.Miner{},
		&miner_repo.Pool{},
		&miner_repo.MinerLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	configFile, err := os.Open("fxhnd.json")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer configFile.Close()
	DevMigrate(postgresDB, configFile)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	routes.RegisterFleetRoutes(postgresDB, router)
	routes.RegisterMinerRoutes(postgresDB, router)
	routes.RegisterScannerRoutes(postgresDB, router)
	routes.RegisterMinerTimeSeriesRoutes(router)

	// Start the router on a separate goroutine
	go func() {
		if err := router.Run(); err != nil {
			log.Printf("Failed to start router: %v", err)
		}
	}()

	// TODO: separate the worker logic from the main function
	//---------------------------WORKER LOGIC--------------------------------

	fleetRepo := fleet_repo.NewFleetRepository(postgresDB)
	// minerRepo := miner_repo.NewMinerRepository(postgresDB)
	// scannerRepo := scanner_repo.NewScannerRepository(postgresDB)
	minerTimeSeriesRepository := miner_repo.NewMinerTimeSeriesRepository(InfluxDBConnectionSettings)

	panicHandler := func(p interface{}) {
		log.Println("worker paniced %v", p)
	}

	pool := pond.New(20, 100, pond.PanicHandler(panicHandler))
	defer pool.StopAndWait()

	workerErrors := make(chan error)
	defer close(workerErrors)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Create a context that can be cancelled
	// ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop() // Ensure that resources are freed on return

	for range ticker.C {
		fmt.Println("Running scheduled tasks...")

		fleets, err := fleetRepo.ListScannersByFleet()
		if err != nil {
			fmt.Println("Error getting fleet list:", err)
			// TODO: notify with alert layer 4
			continue
		}

		for _, fleet := range fleets {

			pool.Submit(func() {
				log.Println("=========================")
				log.Printf("Processing scanner ID: %d\n", fleet.ID)
				log.Println("fleet start ip", fleet.Scanner.Scanner.StartIP)
				log.Println("fleet end ip", fleet.Scanner.Scanner.EndIP)

				startIP := net.ParseIP(fleet.Scanner.Scanner.StartIP)
				endIP := net.ParseIP(fleet.Scanner.Scanner.EndIP)

				if startIP == nil || endIP == nil {
					workerErrors <- fmt.Errorf("invalid IP address format")
					return
				}

				var ips []net.IP
				for ip := startIP.To16(); ; inc(ip) {
					newIP := make(net.IP, len(ip))
					copy(newIP, ip)
					ips = append(ips, newIP)

					if ip.Equal(endIP) {
						break
					}
				}

				log.Println("ips", ips)
				log.Println("ARP request began...")

				// 1, ARP request
				ARPResponses := make(chan *ant_miner_cgi_service.AntminerCGI, len(ips))
				var wg sync.WaitGroup
				for i, ip := range ips {
					wg.Add(1)

					go func(i int, ip net.IP) {
						// 2, make a channel for each miner, which can be reused by the pipeline pattern
						defer wg.Done()
						log.Println("ARP request for", ip, i)

						t := http_auth.NewTransport(fleet.Scanner.Config.Username, fleet.Scanner.Config.Password)
						newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_system_info.cgi", ip), nil)
						if err != nil {
							log.Println("Error creating new request", err)
							return
						}
						resp, err := t.RoundTrip(newRequest)
						if err != nil {
							log.Println("Error in round trip", err)
							return
						}
						defer resp.Body.Close()

						body, err := io.ReadAll(resp.Body)
						if err != nil {
							log.Println("Error reading response body", err)
							return
						}

						var rawGetSystemInfoResponse ant_miner_cgi_queries.RawGetSystemInfoResponse
						if err := json.Unmarshal(body, &rawGetSystemInfoResponse); err != nil {
							log.Println("Error unmarshalling response body", err)
							return
						}

						fmt.Println("rawGetSystemInfoResponse: ", rawGetSystemInfoResponse.MinerType)

						antMinerCGI := ant_miner_cgi_service.NewAntminerCGI(
							miner_domain.Config{
								Username: fleet.Scanner.Config.Username,
								Password: fleet.Scanner.Config.Password,
								Firmware: rawGetSystemInfoResponse.FirmwareType,
							},
							miner_domain.Miner{
								IPAddress:  rawGetSystemInfoResponse.IPAddress,
								MacAddress: rawGetSystemInfoResponse.MACAddress,
							},
							rawGetSystemInfoResponse.MinerType,
						)

						fmt.Println("MODEL NAME =====>>>>>", rawGetSystemInfoResponse.MinerType)
						// feed the ARPResponses channel with the antMinerCGI object
						ARPResponses <- antMinerCGI
					}(i, ip)
				} // end of the ARP request

				wg.Wait()
				close(ARPResponses)

				// 3, Execute all the service codes to update the data fields in the ARPResponses channel
				antMinerCGIServiceChannel := make(chan *ant_miner_cgi_service.AntminerCGI, len(ARPResponses))

				var wgAntMinerCGIService sync.WaitGroup

				for antMinerCGI := range ARPResponses {
					fmt.Println("ant miner service response before executing the commands", *antMinerCGI)
					wgAntMinerCGIService.Add(1)
					go func(antMinerCGI *ant_miner_cgi_service.AntminerCGI) {
						defer wgAntMinerCGIService.Done()

						err := antMinerCGI.CheckStats()
						if err != nil {
							log.Printf("Error checking stats: %v", err)

							workerErrors <- err
							return
						}

						err = antMinerCGI.CheckPools()
						if err != nil {
							log.Printf("Error checking pools: %v", err)

							workerErrors <- err
							return
						}

						err = antMinerCGI.CheckConfig()
						if err != nil {
							log.Printf("Error checking config: %v", err)

							workerErrors <- err
							return
						}

						antMinerCGIServiceChannel <- antMinerCGI

					}(antMinerCGI)
				}

				wgAntMinerCGIService.Wait()
				close(antMinerCGIServiceChannel)

				// make the array for the service channel
				antMinerCGIServiceArray := make([]*ant_miner_cgi_service.AntminerCGI, 0)
				for antMinerCGI := range antMinerCGIServiceChannel {
					antMinerCGIServiceArray = append(antMinerCGIServiceArray, antMinerCGI)
				}

				// 4, Check Alerts and Update the Database
				conditionCounter := make(map[scanner_domain.AlertConditionType]int) // machine count
				// NOTE: if conditions are met, set the alert flag to true. Otherwise, set it to false
				alertFlag := true

				for _, alertCondition := range fleet.Scanner.Alert.Condition {
					conditionCounter[alertCondition.ConditionType] = 0 // default TriggerValue
				}

				for _, antMinerCGIService := range antMinerCGIServiceArray {
					fmt.Println("ant miner alert response before checking the conditions")
					for _, alertCondition := range fleet.Scanner.Alert.Condition {

						if antMinerCGIService.Mode == miner_domain.SleepMode {

							fmt.Println("Skipping the alert service logic with mode", antMinerCGIService.Mode)

							continue
						}

						switch alertCondition.ConditionType {

						case scanner_domain.Hashrate:
							if antMinerCGIService.Stats.HashRate >= float64(alertCondition.TriggerValue) {
								//
								// increment the counter and update the status of miner
								conditionCounter[scanner_domain.Hashrate]++
								antMinerCGIService.Status = miner_domain.HashrateError
							}

						case scanner_domain.Temperature:
							maxTemperature := 0
							for _, temperatureSensor := range antMinerCGIService.Temperature {
								for _, pcbSensor := range temperatureSensor.PcbSensors {
									if pcbSensor >= maxTemperature {
										maxTemperature = pcbSensor
									}
								}
							}

							if maxTemperature >= int(alertCondition.TriggerValue) {
								conditionCounter[scanner_domain.Temperature]++
								antMinerCGIService.Status = miner_domain.TemperatureError
							}

						case scanner_domain.FanSpeed:
							maxFanSpeed := 0
							for _, fanSensor := range antMinerCGIService.Fan {
								if fanSensor.Speed >= maxFanSpeed {
									maxFanSpeed = fanSensor.Speed
								}
							}

							if maxFanSpeed >= int(alertCondition.TriggerValue) {
								conditionCounter[scanner_domain.FanSpeed]++
								antMinerCGIService.Status = miner_domain.FanSpeedError
							}

						case scanner_domain.PoolShares:
							// NOTE: retrieve only the first pool for now, assuming that the pool switch's occureance is rare
							if antMinerCGIService.Pools[0].Accepted <= int(alertCondition.TriggerValue) {
								conditionCounter[scanner_domain.PoolShares]++
								antMinerCGIService.Status = miner_domain.PoolShareError
							}
						}
					}

				}

				for _, alertCondition := range fleet.Scanner.Alert.Condition {
					if conditionCounter[alertCondition.ConditionType] >= int(alertCondition.MachineCount) {
						alertFlag = true
					} else {
						alertFlag = false
					}
				}

				// (4, resolve the alert triggers)
				if alertFlag {
					fleet.Scanner.Alert.State = scanner_domain.Triggered

					result := postgresDB.Where("name = ?", fleet.Scanner.Alert.Name).First(&fleet.Scanner.Alert)
					if result.RowsAffected == 0 {
						err := postgresDB.Create(&fleet.Scanner).Error
						fmt.Println("Error in database create operation", err)
					} else {
						err := postgresDB.Where("name = ?", fleet.Scanner.Alert.Name).Save(&fleet.Scanner.Alert).Error
						fmt.Println("Error in database save operation", err)
					}

					actionCommand := fleet.Scanner.Alert.Action

					// need a way to have a lastUpdatedAt
					var wgAlert sync.WaitGroup
					switch actionCommand {

					// TODO: goroutines for each miner operation
					case scanner_domain.Reboot:

						for _, antMinerCGIService := range antMinerCGIServiceArray {
							wgAlert.Add(1)
							go func(service *ant_miner_cgi_service.AntminerCGI) {
								defer wgAlert.Done()
								err := service.Reboot()
								if err != nil {
									log.Printf("Error rebooting the miner: %v", err)
								}
							}(antMinerCGIService)
						}

					case scanner_domain.Sleep:

						for _, antMinerCGIService := range antMinerCGIServiceArray {
							wgAlert.Add(1)
							go func(service *ant_miner_cgi_service.AntminerCGI) {
								defer wgAlert.Done()
								err := service.SetSleepMode()
								if err != nil {
									log.Printf("Error sleeping the miner: %v", err)
								}
							}(antMinerCGIService)
						}

					case scanner_domain.Normal:

						for _, antMinerCGIService := range antMinerCGIServiceArray {
							wgAlert.Add(1)
							go func(service *ant_miner_cgi_service.AntminerCGI) {
								defer wgAlert.Done()
								err := service.SetNormalMode()
								if err != nil {
									log.Printf("Error setting the miner to normal mode: %v", err)
								}
							}(antMinerCGIService)
						}

					}
					wgAlert.Wait()
				} // end of the case for alertFlag = true

				for _, antMinerCGIService := range antMinerCGIServiceArray {

					var miner miner_repo.Miner
					result := postgresDB.First(&miner, "mac_address = ?", antMinerCGIService.Miner.MacAddress)

					if result.RowsAffected == 0 {
						miner.Miner = miner_domain.Miner{
							MacAddress: antMinerCGIService.Miner.MacAddress,
							IPAddress:  antMinerCGIService.Miner.IPAddress,
						}
						miner.Stats = miner_domain.Stats{
							HashRate:  antMinerCGIService.Stats.HashRate,
							RateIdeal: antMinerCGIService.Stats.RateIdeal,
							Uptime:    antMinerCGIService.Stats.Uptime,
						}
						miner.Config = miner_domain.Config{
							Username: antMinerCGIService.Config.Username,
							Password: antMinerCGIService.Config.Password,
							Firmware: antMinerCGIService.Config.Firmware,
						}

						miner.ModelName = antMinerCGIService.Model
						miner.Mode = miner_domain.NormalMode

						miner.Status = miner_domain.Online
						miner.Pools = []miner_repo.Pool{}
						miner.FleetID = fleet.ID

						// TODO: redo the pool
						for _, pool := range antMinerCGIService.Pools {
							newPool := miner_repo.Pool{
								Pool: miner_domain.Pool{
									Url:      pool.Url,
									User:     pool.User,
									Pass:     pool.Pass,
									Status:   pool.Status,
									Accepted: pool.Accepted,
									Rejected: pool.Rejected,
									Stale:    pool.Stale,
								},
							}

							miner.Pools = append(miner.Pools, newPool)
						}

						// TODO: retrieve max temp and max fan speed
						for _, temperatureSensor := range antMinerCGIService.Temperature {
							miner.Temperature = append(miner.Temperature, temperatureSensor.PcbSensors...)
						}

						for _, fanSensor := range antMinerCGIService.Fan {
							miner.Fan = append(miner.Fan, fanSensor.Speed)
						}

						err := postgresDB.Create(&miner).Error
						if err != nil {
							fmt.Println("error in miner create operation")
							workerErrors <- err
						}

						minerTimeSeriesRepository.WriteMinerData(miner.Miner.MacAddress, miner_repo.MinerTimeSeries{
							MacAddress: antMinerCGIService.Miner.MacAddress,
							HashRate:   int(antMinerCGIService.Stats.HashRate),
							TempSensor: miner.Temperature,
							FanSensor:  miner.Fan,
						})

						if antMinerCGIService.Pools != nil {
							minerTimeSeriesRepository.WritePoolData(miner.Miner.MacAddress, miner_repo.PoolTimeSeries{
								MacAddress: antMinerCGIService.Miner.MacAddress,
								Accepted:   antMinerCGIService.Pools[0].Accepted,
								Rejected:   antMinerCGIService.Pools[0].Rejected,
								Stale:      antMinerCGIService.Pools[0].Stale,
							})

						} else {
							fmt.Println("No pool data available")
							minerTimeSeriesRepository.WritePoolData(miner.Miner.MacAddress, miner_repo.PoolTimeSeries{
								MacAddress: antMinerCGIService.Miner.MacAddress,
								Accepted:   0,
								Rejected:   0,
								Stale:      0,
							})
						}

						// result.RowsAffected != 0
						// = a relevant miner already exists
					} else {

						var existingMiner miner_repo.Miner

						err := postgresDB.Model(&miner_repo.Miner{}).
							Preload("Pools").
							Find(&existingMiner, miner.ID).Error
						if err != nil {
							fmt.Println("error in preload models ", err)
						}

						existingMiner.Miner.IPAddress = antMinerCGIService.Miner.IPAddress
						existingMiner.Miner.MacAddress = antMinerCGIService.Miner.MacAddress
						existingMiner.Stats.HashRate = antMinerCGIService.Stats.HashRate
						existingMiner.Stats.RateIdeal = antMinerCGIService.Stats.RateIdeal
						existingMiner.Stats.Uptime = antMinerCGIService.Stats.Uptime
						existingMiner.Config.Username = antMinerCGIService.Config.Username
						existingMiner.Config.Password = antMinerCGIService.Config.Password
						existingMiner.Config.Firmware = antMinerCGIService.Config.Firmware
						existingMiner.Mode = antMinerCGIService.Mode
						existingMiner.ModelName = antMinerCGIService.Model

						existingMiner.Status = antMinerCGIService.Status
						existingMiner.FleetID = fleet.ID

						for index, pool := range antMinerCGIService.Pools {
							existingMiner.Pools[index].Pool.Url = pool.Url
							existingMiner.Pools[index].Pool.User = pool.User
							existingMiner.Pools[index].Pool.Pass = pool.Pass
							existingMiner.Pools[index].Pool.Status = pool.Status
							existingMiner.Pools[index].Pool.Accepted = pool.Accepted
							existingMiner.Pools[index].Pool.Rejected = pool.Rejected
							existingMiner.Pools[index].Pool.Stale = pool.Stale

							postgresDB.Where("miner_id = ?", existingMiner.ID).Save(existingMiner.Pools[index])
						}

						for _, temperatureSensor := range antMinerCGIService.Temperature {
							// TODO: max search
							existingMiner.Temperature = append(miner.Temperature, temperatureSensor.PcbSensors...)
						}

						for _, fanSensor := range antMinerCGIService.Fan {
							// TODO: max search
							existingMiner.Fan = append(miner.Fan, fanSensor.Speed)
						}

						if err := postgresDB.Where("ID = ?", existingMiner.ID).Updates(existingMiner).Error; err != nil {
							fmt.Println("error in seesssion ", err)
						}

						// update the timeseries data for the existing miner
						minerTimeSeriesRepository.WriteMinerData(existingMiner.Miner.MacAddress, miner_repo.MinerTimeSeries{
							MacAddress: antMinerCGIService.Miner.MacAddress,
							HashRate:   int(antMinerCGIService.Stats.HashRate),
							TempSensor: existingMiner.Temperature, // reuse the existing temperature data
							FanSensor:  existingMiner.Fan,         // reuse the existing fan data
						})

						if antMinerCGIService.Pools != nil {

							minerTimeSeriesRepository.WritePoolData(miner.Miner.MacAddress, miner_repo.PoolTimeSeries{
								MacAddress: antMinerCGIService.Miner.MacAddress,
								Accepted:   antMinerCGIService.Pools[0].Accepted,
								Rejected:   antMinerCGIService.Pools[0].Rejected,
								Stale:      antMinerCGIService.Pools[0].Stale,
							})

						} else {

							minerTimeSeriesRepository.WritePoolData(miner.Miner.MacAddress, miner_repo.PoolTimeSeries{
								MacAddress: antMinerCGIService.Miner.MacAddress,
								Accepted:   0,
								Rejected:   0,
								Stale:      0,
							})

						}
					}
				}
				// Flush the timeseries data to the database
				minerTimeSeriesRepository.FlushMinerData()
				minerTimeSeriesRepository.FlushPoolData()

				fmt.Println("========================END OF WORKER=========================")
			})
		}
	}

}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
