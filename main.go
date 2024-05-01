package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/alitto/pond"
	"github.com/gin-gonic/gin"

	postgres "foxhound/internal/infrastructure/database/postgres"
	http_auth "foxhound/pkg/http_auth"

	// TODO: db migration/seed
	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	miner_repo "foxhound/internal/infrastructure/database/repositories/miner"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"

	miner_domain "foxhound/internal/application/miner/domain"
	scanner_domain "foxhound/internal/application/scanner/domain"

	ant_miner_cgi_queries "foxhound/internal/application/miner/ant_miner_cgi/queries"
	ant_miner_cgi_service "foxhound/internal/application/miner/ant_miner_cgi/service"
)

// TODO: select statement for different vendors
// TODO: R&D for pool library's memory leak
// TODO: logic for identifying the active pool
// TODO: logic for pool change

func main() {

	postgresDB := postgres.Init()
	err := postgresDB.AutoMigrate(
		// NOTE: The order matters
		&fleet_repo.Fleet{},

		&scanner_repo.Scanner{},
		&scanner_repo.Alert{},
		&scanner_repo.AlertCondition{},
		&scanner_repo.AlertLog{},

		&miner_repo.Miner{},
		&miner_repo.Pool{},
		&miner_repo.TemperatureSensor{},
		&miner_repo.PcbSensor{},
		&miner_repo.FanSensor{},
		&miner_repo.MinerLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	DevMigrate(postgresDB)

	router := gin.Default()

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	// Create a context that can be cancelled
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel() // Ensure that resources are freed on return

	fleetRepo := fleet_repo.NewFleetRepository(postgresDB)
	// minerRepo := miner_repo.NewMinerRepository(postgresDB)
	// scannerRepo := scanner_repo.NewScannerRepository(postgresDB)

	pool := pond.New(20, 100)
	defer pool.StopAndWait()

	workerErrors := make(chan error)

	// ticker := time.NewTicker(20 * time.Second)
	// defer ticker.Stop()
	// for {
	// 	select {
	// 	// TODO: cancel logic for user activity
	// 	case <-ctx.Done():
	// 		fmt.Println("Shutting down the scheduled tasks...")
	// 		return

	// 	case <-ticker.C:
	// 		fmt.Println("Running scheduled tasks...")
	fleets, err := fleetRepo.List()
	if err != nil {
		// alert layer =  4
		fmt.Println("Error getting fleet list:", err)
		// continue
	}

	fmt.Println("fleet list", fleets)

	for _, fleet := range fleets {

		pool.Submit(func() {
			log.Println("=========================")
			log.Printf("Processing scanner ID: %d\n", fleet.ID)
			log.Println("fleet", fleet)

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

					fmt.Println("rawGetSystemInfoResponse: ", rawGetSystemInfoResponse)

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
					)
					// feed the ARPResponses channel with the antMinerCGI object
					ARPResponses <- antMinerCGI
				}(i, ip)
			} // end of the ARP request

			wg.Wait()
			close(ARPResponses)

			// 3, Execute all the service codes to update the data fields in the ARPResponses channel
			// TODO: make two channels : one for alert and the other one for updating the database
			antMinerCGIServiceChannel := make(chan *ant_miner_cgi_service.AntminerCGI, len(ARPResponses))
			var wgAntMinerCGIService sync.WaitGroup

			for antMinerCGIService := range ARPResponses {
				fmt.Println("ant miner service response before executing the commands", *antMinerCGIService)
				wgAntMinerCGIService.Add(1)
				go func(antMinerCGIService *ant_miner_cgi_service.AntminerCGI) {
					defer wgAntMinerCGIService.Done()

					err := antMinerCGIService.CheckStats()
					if err != nil {
						log.Printf("Error checking stats: %v", err)

						workerErrors <- err
						return
					}
					err = antMinerCGIService.CheckPools()
					if err != nil {
						log.Printf("Error checking pools: %v", err)

						workerErrors <- err
						return
					}
					err = antMinerCGIService.CheckConfig()
					if err != nil {
						log.Printf("Error checking config: %v", err)

						workerErrors <- err
						return
					}

					log.Println("antMinerCGIService hashrate:", antMinerCGIService.Stats.HashRate)
					log.Println("antMinerCGIService pool:", antMinerCGIService.Pools[0].Url)
					log.Println("antMinerCGIService username:", antMinerCGIService.Config.Username)
					log.Println("antMInerCGIService firmware", antMinerCGIService.Config.Firmware)

					log.Printf("Submitting miner info to channel: %+v", antMinerCGIService)
					antMinerCGIServiceChannel <- antMinerCGIService

				}(antMinerCGIService)
			}
			// waiting for the ongoing antMinerCGIService Channel
			wgAntMinerCGIService.Wait()
			close(antMinerCGIServiceChannel)

			log.Println("------MINER INFO-------")
			for minerInfo := range antMinerCGIServiceChannel {

				log.Println(minerInfo)

				log.Println("=======INFO END =======")

			}

			// 4, Check Alerts and Update the Database
			conditionCounter := make(map[scanner_domain.AlertConditionType]int) // machine count
			// NOTE: if conditions are met, set the alert flags to true. Otherwise, set it to false
			alertFlag := true

			for _, alertCondition := range fleet.Scanner.Alert.Condition {
				conditionCounter[alertCondition.ConditionType] = 0 // default TriggerValue
			}

			for antMinerCGIService := range antMinerCGIServiceChannel {

				fmt.Println("ant miner alert response before checking the conditions")
				for _, alertCondition := range fleet.Scanner.Alert.Condition {
					switch alertCondition.ConditionType {

					case scanner_domain.Hashrate:

						if antMinerCGIService.Stats.HashRate >= float64(alertCondition.TriggerValue) {
							conditionCounter[scanner_domain.Hashrate]++
						}

					case scanner_domain.Temperature:

						maxTemperature := 0
						for _, temperatureSensor := range antMinerCGIService.Temperature {
							for _, pcbSensor := range temperatureSensor.PcbSensors {
								if pcbSensor.Temperature > maxTemperature {
									maxTemperature = pcbSensor.Temperature
								}
							}
						}

						if maxTemperature >= int(alertCondition.TriggerValue) {
							conditionCounter[scanner_domain.Temperature]++
						}

					case scanner_domain.FanSpeed:

						maxFanSpeed := 0
						for _, fanSensor := range antMinerCGIService.Fan {
							if fanSensor.Speed > maxFanSpeed {
								maxFanSpeed = fanSensor.Speed
							}
						}

						if maxFanSpeed >= int(alertCondition.TriggerValue) {
							conditionCounter[scanner_domain.FanSpeed]++
						}

					case scanner_domain.PoolShares:
						// NOTE: retrieve only the first pool for now, assuming that the pool switch's occureance is rare
						if antMinerCGIService.Pools[0].Accepted <= int(alertCondition.TriggerValue) {
							conditionCounter[scanner_domain.PoolShares]++
						}
						// case scanner_domain.OfflineMiners:
					}
				}

				//  feeding the updated fleet payload for updating the database
				fleet.Miners = append(fleet.Miners, miner_repo.Miner{
					Miner:  antMinerCGIService.Miner,
					Stats:  antMinerCGIService.Stats,
					Config: antMinerCGIService.Config,
					Mode:   antMinerCGIService.Mode, // TODO: update logic should be done in the service method in the interface
					Status: antMinerCGIService.Status,
				})
			}

			for _, alertCondition := range fleet.Scanner.Alert.Condition {
				if conditionCounter[alertCondition.ConditionType] >= int(alertCondition.MachineCount) {
					// if all the conditions are met, the alertFlag
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

				switch actionCommand {

				case scanner_domain.Reboot:

				case scanner_domain.Sleep:

				case scanner_domain.Normal:

				case scanner_domain.ChangePool:

				}

			}

			// 5, Update the database
			//     the fleet stats (after the triggers get activated) gets reflected in the next periodic update
			// TODO: batch - tx. operation
			for _, miner := range fleet.Miners {
				result := postgresDB.Where("mac_address = ?", miner.Miner.MacAddress).First(&miner)

				if result.RowsAffected == 0 {
					err := postgresDB.Create(&miner).Error
					fmt.Println("Error in database create operation", err)
				} else {
					err := postgresDB.Where("mac_address = ?", miner.Miner.MacAddress).Save(&miner).Error
					fmt.Println("Error in database save operation", err)
				}

			}

		})
	}

	// After updating the miner payload (in the application layer),
	// update the miner payload to the database with upsert operation

	// kill go routines after the fleet -> scanner -> miner list is processed
}

// }
// }

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
