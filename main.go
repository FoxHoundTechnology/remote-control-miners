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
					// fill the struct inside of the channel
					ARPResponses <- antMinerCGI
				}(i, ip)
			} // end of the ARP request

			wg.Wait()
			close(ARPResponses)

			// 3, Execute all the service codes to update the data fields in the ARPResponses channel
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

			// 4, Check Alerts
			// create a map object that covers alert conditions with flag value
			// antMinerCGIAlertsChannel := make(chan *ant_miner_cgi_service.AntminerCGI)
			conditionCounter := make([]scanner_repo.AlertCondition, 0, len(fleet.Scanner.Alert.Condition))
			for _, alertCondition := range fleet.Scanner.Alert.Condition {
				conditionCounter = append(conditionCounter, alertCondition)
			}
			alertFlags := make(map[scanner_domain.AlertConditionType]bool, len(conditionCounter))
			for _, alertCondition := range conditionCounter {
				alertFlags[alertCondition.Condition] = false
			}

			for antMinerCGIService := range antMinerCGIServiceChannel {
				fmt.Println("ant miner alert response before checking the conditions")
				for index, alertCondition := range fleet.Scanner.Alert.Condition {
					switch alertCondition.Condition {

					case scanner_domain.Hashrate:
						if antMinerCGIService.Stats.HashRate > 100 {
							conditionCounter[index].Value = conditionCounter[index].Value + 1
						}

					case scanner_domain.Temperature:
						// for _, temperature := range antminercgiservice.temperature {

						// 	// find the highest temp in each chain
						// 	// and if the highest temp is greater than the alert condition value
						// 	// increment the condition value

						// }
					case scanner_domain.FanSpeed:
						// for _, fan := range antMinerCGIService.Fan {
						// 	if fan.Value > alertCondition.Value {

						// 	}
						// }
					case scanner_domain.PoolShares:
						// for _, pool := range antMinerCGIService.Pools {
						// if pool.Shares < alertCondition.Value {
						// 	alertFlag = true
						// }
						// }
					}
				}

			}

			// wg.Wait()
			// close(antMinerCGIAlertsChannel)

			// 4, Update the database

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
