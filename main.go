package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/gin-gonic/gin"

	postgres "foxhound/internal/infrastructure/database/postgres"

	"foxhound/pkg/http_auth"

	// TODO: db migration/seed
	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	miner_repo "foxhound/internal/infrastructure/database/repositories/miner"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"

	miner_domain "foxhound/internal/application/miner/domain"

	ant_miner_cgi_queries "foxhound/internal/application/miner/ant_miner_cgi/queries"
	ant_miner_cgi_service "foxhound/internal/application/miner/ant_miner_cgi/service"
)

// TODO: R&D for pool library's memory leak
func main() {

	postgresDB := postgres.Init()

	err := postgresDB.AutoMigrate(
		// NOTE: The order matters
		&fleet_repo.Fleet{},

		&scanner_repo.Scanner{},
		&scanner_repo.Alert{},
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure that resources are freed on return

	fleetRepo := fleet_repo.NewFleetRepository(postgresDB)
	// minerRepo := miner_repo.NewMinerRepository(postgresDB)
	// scannerRepo := scanner_repo.NewScannerRepository(postgresDB)

	pool := pond.New(10, 1000)
	workerErrors := make(chan error)
	// go func() { }()
	// ticker
	ticker := time.NewTicker(80 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down the scheduled tasks...")
			return

		case <-ticker.C:
			fmt.Println("Running scheduled tasks...")
			fleets, err := fleetRepo.List()
			if err != nil {
				// alert layer =  4
				fmt.Println("Error getting fleet list:", err)
				continue
			}

			fmt.Println("fleet list", fleets)
			// spawn a goroutine for each fleet -> scanner
			for _, fleet := range fleets {
				time.Sleep(30 * time.Second)
				pool.Submit(func() {
					fmt.Println("=========================")
					// TODO: select statement for different vendors

					// initialize the antMinerController
					fmt.Printf("Processing scanner ID: %d\n", fleet.ID)
					fmt.Println("fleet", fleet)

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

					fmt.Println("ips", ips)
					fmt.Println("ARP request began...")

					// 1, ARP request
					ARPResponses := make(chan *ant_miner_cgi_service.AntminerCGI, len(ips))
					var wg sync.WaitGroup
					for i, ip := range ips {
						wg.Add(1)
						time.Sleep(1 * time.Second)

						go func(i int, ip net.IP) {
							// 2, make channel for each miner that keeps get reused by the pipeline pattern
							defer wg.Done()
							fmt.Println("ARP request for", ip, i)

							// get system info
							t := http_auth.NewTransport(fleet.Scanner.Config.Username, fleet.Scanner.Config.Password)
							newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_system_info.cgi", ip), nil)
							if err != nil {
								log.Println("Error creating new request", err)
							}
							resp, err := t.RoundTrip(newRequest)
							if err != nil {
								log.Println("Error in round trip", err)
							}
							defer resp.Body.Close()

							body, err := io.ReadAll(resp.Body)
							if err != nil {
								log.Println("Error reading response body", err)
							}

							var rawGetSystemInfoResponse ant_miner_cgi_queries.RawGetSystemInfoResponse
							if err := json.Unmarshal(body, &rawGetSystemInfoResponse); err != nil {
								log.Println("Error unmarshalling response body", err)
							}

							fmt.Println("rawGetSystemInfoResponse: ", rawGetSystemInfoResponse)

							// fill the struct inside of the channel
							ARPResponses <- &ant_miner_cgi_service.AntminerCGI{
								Miner: miner_domain.Miner{
									IPAddress:  rawGetSystemInfoResponse.IPAddress,
									MacAddress: rawGetSystemInfoResponse.MACAddress,
								},
								Config: miner_domain.Config{
									Username: fleet.Scanner.Config.Username,
									Password: fleet.Scanner.Config.Password,
									Firmware: rawGetSystemInfoResponse.FirmwareType,
								},
							}

						}(i, ip)
					} // end of the ARP request

					wg.Wait()
					close(ARPResponses)

					// 2, execute the all the service code to update the data fields in the ARPResponses channel
					antMinerCGIServiceChannel := make(chan *ant_miner_cgi_service.AntminerCGI)
					for antMinerCGIService := range ARPResponses {
						fmt.Println("ant miner service response before executing the commands", *antMinerCGIService)
						wg.Add(1)
						go func(antMinerCGIService *ant_miner_cgi_service.AntminerCGI) {
							defer wg.Done()

							err := antMinerCGIService.CheckStats()
							if err != nil {
								workerErrors <- err
							}
							err = antMinerCGIService.CheckPools()
							if err != nil {
								workerErrors <- err
							}
							err = antMinerCGIService.CheckConfig()
							if err != nil {
								workerErrors <- err
							}

							fmt.Println("antMinerCGIService hashrate:", antMinerCGIService.Stats.HashRate)
							fmt.Println("antMinerCGIService pool:", antMinerCGIService.Pools[0].Url)
							fmt.Println("antMinerCGIService username:", antMinerCGIService.Config.Username)
							fmt.Println("updating the data in the ARP channel and aggregate them all into another channel: antMinerCGIServiceCHannel ")

							antMinerCGIServiceChannel <- antMinerCGIService

						}(antMinerCGIService)
					}
					// waiting for the ongoing antMinerCGIService Channel
					wg.Wait()
					close(antMinerCGIServiceChannel)

					// 3, Check Alerts
					// create a map object that covers condition map data with flag value

					antMinerCGIAlertsChannel := make(chan *ant_miner_cgi_service.AntminerCGI)
					for antMinerCGIAlert := range antMinerCGIServiceChannel {
						fmt.Println("ant miner alert response before checking the conditions ")
						wg.Add(1)

						go func(antMinerCGIAlert *ant_miner_cgi_service.AntminerCGI) {
							defer wg.Done()
							// for loop through the preconfigured alerts condition

						}(antMinerCGIAlert)

					}

					wg.Wait()
					close(antMinerCGIAlertsChannel)

					// 4, update the database

				})
			}

			// After updating the miner payload (in the application layer),
			// update the miner payload to the database with upsert operation

			// kill go routines after the fleet -> scanner -> miner list is processed
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
