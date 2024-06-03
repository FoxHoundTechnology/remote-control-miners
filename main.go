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

	timeseries "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/influxdb"
	postgres "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/postgres"

	http_auth "github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	// TODO: db migration/seed
	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	scanner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/scanner/domain"

	fleet_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/fleet"
	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"
	scanner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/scanner"

	routes "github.com/FoxHoundTechnology/remote-control-miners/internal/interface/routers"

	ant_miner_cgi_queries "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/queries"
	ant_miner_cgi_service "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/service"
)

// TODO: retrieve a collection of max temperature and max fan speed
// TODO: select statement for different vendors
// TODO: batch operation for miner stats update
// TODO: alert layer support
// TODO: R&D for pool library's memory leak
// TODO: create a list object for pool function calls
// TODO: logic for detecting offline miners
// TODO: logic for identifying the active pool
// TODO: logic for combined miner error supports

func main() {

	postgresDB := postgres.Init()
	InfluxDBConnectionSettings := timeseries.Init()

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

	// dbConnectionSetting, err := postgresDB.DB()
	// if err != nil {
	// 	log.Fatalf("Failed to get the database connection: %v", err)
	// }

	// dbConnectionSetting.SetMaxIdleConns(30)
	// dbConnectionSetting.SetMaxOpenConns(60)
	// dbConnectionSetting.SetConnMaxLifetime(time.Hour)

	// // TODO: refactor unique contraint logic
	// err = postgresDB.Exec(miner_repo.CreateUniqueMinerIndexSQL).Error
	// if err != nil {
	// 	fmt.Println("Error creating unique index for miners", err)
	// }

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

	go func() {
		if err := router.Run(); err != nil {
			log.Printf("Failed to start router: %v", err)
		}
	}()

	// TODO: separate the worker logic from the main function
	//---------------------------WORKER LOGIC--------------------------------
	fleetRepo := fleet_repo.NewFleetRepository(postgresDB)
	minerTimeSeriesRepository := miner_repo.NewMinerTimeSeriesRepository(InfluxDBConnectionSettings)

	panicHandler := func(p interface{}) {
		log.Println("worker paniced %v", p)
	}

	pool := pond.New(
		200,
		100,
		pond.PanicHandler(panicHandler),
		pond.Strategy(pond.Lazy()),
		pond.MinWorkers(29),
	)
	defer pool.StopAndWait()

	workerErrors := make(chan error)
	defer close(workerErrors)

	ticker := time.NewTicker(300 * time.Second)
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
			fleet := fleet
			pool.Submit(func() {
				log.Println("===========current pool worker #==============", pool.RunningWorkers())

				fmt.Println("Worker Pool Metrics:")
				fmt.Printf("Running Workers: %d\n", pool.RunningWorkers())
				fmt.Printf("Idle Workers: %d\n", pool.IdleWorkers())
				fmt.Printf("Min Workers: %d\n", pool.MinWorkers())
				fmt.Printf("Max Workers: %d\n", pool.MaxWorkers())
				fmt.Printf("Max Capacity: %d\n", pool.MaxCapacity())
				fmt.Printf("Submitted Tasks: %d\n", pool.SubmittedTasks())
				fmt.Printf("Waiting Tasks: %d\n", pool.WaitingTasks())
				fmt.Printf("Successful Tasks: %d\n", pool.SuccessfulTasks())
				fmt.Printf("Failed Tasks: %d\n", pool.FailedTasks())
				fmt.Printf("Completed Tasks: %d\n", pool.CompletedTasks())

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

				ARPResponses := make(chan *ant_miner_cgi_service.AntminerCGI, len(ips))
				var wg sync.WaitGroup
				for i, ip := range ips {
					wg.Add(1)
					go func(i int, ip net.IP) {
						defer wg.Done()

						t := http_auth.NewTransport(fleet.Scanner.Config.Username, fleet.Scanner.Config.Password)
						newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_system_info.cgi", ip), nil)
						if err != nil {
							log.Println("Error creating new request", err)
							return
						}
						resp, err := t.RoundTrip(newRequest)
						if err != nil {
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
					antMinerCGIService := antMinerCGIService

					for _, alertCondition := range fleet.Scanner.Alert.Condition {

						if antMinerCGIService.Mode == miner_domain.SleepMode {
							continue
						}

						switch alertCondition.ConditionType {

						case scanner_domain.Hashrate:
							if antMinerCGIService.Stats.HashRate <= float64(alertCondition.TriggerValue) {
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

				fmt.Println("alertFlag", alertFlag)

				/// TODO: alert layer support
				if false {
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

				minerRepository := miner_repo.NewMinerRepository(postgresDB)
				existingMiners, err := minerRepository.ListByFleetID(fleet.ID)
				if err != nil {
					fmt.Println("Error getting miner list:", err)
					workerErrors <- err
				}

				var updatedMiners []*miner_repo.Miner
				var newMiners []*miner_repo.Miner

				fmt.Println("-------------- fleet ", fleet.ID, "------------------")
				fmt.Println("Miners in the fleet", existingMiners)
				fmt.Println("----------------------------------------------------")

				// get the list of miners based on the fleet ID
				for _, updatedMiner := range antMinerCGIServiceArray {
					found := false
					// base case
					if len(existingMiners) != 0 {
						fmt.Println("Initial scanning, no miners found in the database")
						for _, existingMiner := range existingMiners {
							if updatedMiner.Miner.MacAddress == existingMiner.Miner.MacAddress {
								// update the miner data
								existingMiner.ModelName = updatedMiner.Model
								existingMiner.Miner.IPAddress = updatedMiner.Miner.IPAddress
								existingMiner.Status = updatedMiner.Status
								existingMiner.Mode = updatedMiner.Mode
								existingMiner.MinerType = 0

								existingMiner.Stats.HashRate = updatedMiner.Stats.HashRate
								existingMiner.Stats.RateIdeal = updatedMiner.Stats.RateIdeal
								existingMiner.FleetID = fleet.ID

								existingMiner.Fan = make([]int, len(updatedMiner.Fan))
								for i, fan := range updatedMiner.Fan {
									existingMiner.Fan[i] = fan.Speed
								}

								existingMiner.Temperature = make([]int, len(updatedMiner.Temperature))
								for i, temperature := range updatedMiner.Temperature {
									max := 0
									for _, pcbSensor := range temperature.PcbSensors {
										if pcbSensor >= max {
											max = pcbSensor
										}
									}
									existingMiner.Temperature[i] = max
								}

								// NOTE: avoid out of index error
								if updatedMiner.Pools != nil {
									existingMiner.Pools = make([]miner_repo.Pool, len(updatedMiner.Pools))
									for i, pool := range updatedMiner.Pools {
										existingMiner.Pools[i].Pool = miner_domain.Pool{
											Url:      pool.Url,
											User:     pool.User,
											Pass:     pool.Pass,
											Status:   pool.Status,
											Accepted: pool.Accepted,
											Rejected: pool.Rejected,
											Stale:    pool.Stale,
										}
									}
								}

								updatedMiners = append(updatedMiners, existingMiner)
								minerTimeSeriesRepository.WriteMinerData(
									miner_repo.MinerTimeSeries{
										MacAddress: existingMiner.Miner.MacAddress,
										HashRate:   int(existingMiner.Stats.HashRate),
										TempSensor: existingMiner.Temperature,
										FanSensor:  existingMiner.Fan,
									},
								)
								minerTimeSeriesRepository.WritePoolData(
									miner_repo.PoolTimeSeries{
										MacAddress: existingMiner.Miner.MacAddress,
										Accepted:   existingMiner.Pools[0].Pool.Accepted,
										Rejected:   existingMiner.Pools[0].Pool.Rejected,
										Stale:      existingMiner.Pools[0].Pool.Stale,
									},
								)

								found = true
								break
							}
						}
					}

					// register the new miner
					if !found {
						newMiner := &miner_repo.Miner{
							Miner: miner_domain.Miner{
								IPAddress:  updatedMiner.Miner.IPAddress,
								MacAddress: updatedMiner.Miner.MacAddress,
							},
							Stats: miner_domain.Stats{
								HashRate:  updatedMiner.Stats.HashRate,
								RateIdeal: updatedMiner.Stats.RateIdeal,
								Uptime:    updatedMiner.Stats.Uptime,
							},
							Config: miner_domain.Config{
								Username: updatedMiner.Config.Username,
								Password: updatedMiner.Config.Password,
								Firmware: updatedMiner.Config.Firmware,
							},
							ModelName: updatedMiner.Model,
							Mode:      updatedMiner.Mode,

							Status:  updatedMiner.Status,
							FleetID: fleet.ID,
						}

						newMiner.Fan = make([]int, len(updatedMiner.Fan))
						for i, fan := range updatedMiner.Fan {
							newMiner.Fan[i] = fan.Speed
						}

						newMiner.Temperature = make([]int, len(updatedMiner.Temperature))
						for i, temperature := range updatedMiner.Temperature {
							max := 0
							for _, pcbSensor := range temperature.PcbSensors {
								if pcbSensor >= max {
									max = pcbSensor
								}
							}
							newMiner.Temperature[i] = max
						}

						if updatedMiner.Pools != nil {
							newMiner.Pools = []miner_repo.Pool{}
							newMiner.Pools = make([]miner_repo.Pool, len(updatedMiner.Pools))
							for i, pool := range updatedMiner.Pools {
								newMiner.Pools[i].Pool = miner_domain.Pool{
									Url:      pool.Url,
									User:     pool.User,
									Pass:     pool.Pass,
									Status:   pool.Status,
									Accepted: pool.Accepted,
									Rejected: pool.Rejected,
									Stale:    pool.Stale,
								}
							}
						}

						newMiners = append(newMiners, newMiner)

						// timeseries updates
						minerTimeSeriesRepository.WriteMinerData(
							miner_repo.MinerTimeSeries{
								MacAddress: newMiner.Miner.MacAddress,
								HashRate:   int(newMiner.Stats.HashRate),
								TempSensor: newMiner.Temperature,
								FanSensor:  newMiner.Fan,
							},
						)

						minerTimeSeriesRepository.WritePoolData(
							miner_repo.PoolTimeSeries{
								MacAddress: newMiner.Miner.MacAddress,
								Accepted:   newMiner.Pools[0].Pool.Accepted,
								Rejected:   newMiner.Pools[0].Pool.Rejected,
								Stale:      newMiner.Pools[0].Pool.Stale,
							},
						)
					}
				}

				fmt.Println("-------------- TIMESERIES DATABASE OPERATION START: fleet ", fleet.ID, "------------------")
				fmt.Println("up dated miners --->>> ", updatedMiners)
				fmt.Println("new miners --->>> ", newMiners)

				go func() {
					minerTimeSeriesRepository.FlushMinerData()
					minerTimeSeriesRepository.FlushPoolData()
				}()

				minerRepository.CreateMinersInBatch(newMiners)
				minerRepository.UpdateMinersInBatch(updatedMiners)

				fmt.Println("========================END OF WORKER=========================", fleet.Name)
			}) // end of pool submit

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
