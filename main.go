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

	// "github.com/alitto/pond"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	postgres "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/postgres"

	http_auth "github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	// TODO: db migration/seed
	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"

	fleet_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/fleet"
	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"
	scanner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/scanner"

	routes "github.com/FoxHoundTechnology/remote-control-miners/internal/interface/routers"

	ant_miner_cgi_queries "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/queries"
	ant_miner_cgi_service "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/service"
)

// TODO: replace queryMiner with a buffered channel for term ui implementation
// TODO: modify the connection setting for db with gorm API

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
	// InfluxDBConnectionSettings := timeseries.Init()
	dbConfig, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get the database connection: %v", err)
	}
	dbConfig.SetMaxIdleConns(50)

	// Set connection pool settings

	err = postgresDB.AutoMigrate(
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
	err = postgresDB.Exec(miner_repo.CreateUniqueMinerIndexSQL).Error
	if err != nil {
		fmt.Println("Error creating unique index for miners", err)
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

	go func() {
		if err := router.Run(); err != nil {
			log.Printf("Failed to start router: %v", err)
		}
	}()

	// TODO: separate the worker logic from the main function
	//---------------------------WORKER LOGIC--------------------------------
	fleetRepo := fleet_repo.NewFleetRepository(postgresDB)
	// minerTimeSeriesRepository := miner_repo.NewMinerTimeSeriesRepository(InfluxDBConnectionSettings)

	// panicHandler := func(p interface{}) {
	// 	log.Println("worker paniced %v", p)
	// }

	// pool := pond.New(
	// 	200,
	// 	100,
	// 	pond.PanicHandler(panicHandler),
	// 	pond.Strategy(pond.Eager()),
	// 	pond.MinWorkers(29),
	// )
	// defer pool.StopAndWait()

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
			time.Sleep(500 * time.Millisecond) // Adjust the duration as needed

			go func(fleetModel fleet_repo.Fleet) {
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

				minerRepository := miner_repo.NewMinerRepository(postgresDB)
				antMinerCGIModel := make(chan *miner_repo.Miner, len(ips))
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

						err = antMinerCGI.CheckStats()
						if err != nil {
							log.Printf("Error checking stats: %v", err)
							workerErrors <- err
							//	return
						}

						err = antMinerCGI.CheckPools()
						if err != nil {
							log.Printf("Error checking pools: %v", err)
							workerErrors <- err
							//return
						}

						err = antMinerCGI.CheckConfig()
						if err != nil {
							log.Printf("Error checking config: %v", err)
							workerErrors <- err
							//return
						}

						newMinerModel := &miner_repo.Miner{
							Miner: miner_domain.Miner{
								IPAddress:  antMinerCGI.Miner.IPAddress,
								MacAddress: antMinerCGI.Miner.MacAddress,
							},
							Stats: miner_domain.Stats{
								HashRate:  antMinerCGI.Stats.HashRate,
								RateIdeal: antMinerCGI.Stats.RateIdeal,
								Uptime:    antMinerCGI.Stats.Uptime,
							},
							Config: miner_domain.Config{
								Username: antMinerCGI.Config.Username,
								Password: antMinerCGI.Config.Password,
								Firmware: antMinerCGI.Config.Firmware,
							},
							ModelName: antMinerCGI.Model,
							Mode:      antMinerCGI.Mode,

							Status:  antMinerCGI.Status,
							FleetID: fleet.ID,
						}

						newMinerModel.Fan = make([]int, len(antMinerCGI.Fan))
						for i, fan := range antMinerCGI.Fan {
							newMinerModel.Fan[i] = fan.Speed
						}

						newMinerModel.Temperature = make([]int, len(antMinerCGI.Temperature))
						for i, temp := range antMinerCGI.Temperature {
							max := 0
							for _, pcbSensor := range temp.PcbSensors {
								if pcbSensor > max {
									max = pcbSensor
								}
							}
							newMinerModel.Temperature[i] = max
						}

						if newMinerModel.Pools != nil {
							newMinerModel.Pools = make([]miner_repo.Pool, len(newMinerModel.Pools))
							for i, pool := range antMinerCGI.Pools {
								newMinerModel.Pools[i].Pool = miner_domain.Pool{
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

						// feed the ARPResponses channel with the antMinerCGI object
						antMinerCGIModel <- newMinerModel
					}(i, ip)
				} // end of the ARP request

				wg.Wait()
				close(antMinerCGIModel)

				log.Println("length is ", len(antMinerCGIModel))

				minerModelArr := make([]*miner_repo.Miner, len(antMinerCGIModel))
				for antMinerCGI := range antMinerCGIModel {
					minerModelArr = append(minerModelArr, antMinerCGI)
				}

				minerRepository.UpdateMinersInBatch(minerModelArr)
				// minerRepository.CreateMinersInBatch(minerModelArr)

				fmt.Println("========================END OF WORKER=========================", fleet.Name)
			}(fleet)

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
