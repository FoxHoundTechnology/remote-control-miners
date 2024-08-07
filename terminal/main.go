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

	tea "github.com/charmbracelet/bubbletea"

	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"

	ant_miner_cgi_queries "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/queries"
	ant_miner_cgi_service "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/service"

	"github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	// antminer cgi

	config "github.com/FoxHoundTechnology/remote-control-miners/terminal/util"
)

var INTERVAL_MINS = 2
var MAX_WORKERS = 10

// It's used to update the model with a new set of miners.
type minerListMsg []*miner_repo.Miner

// It doesn't carry any data itself; it's just a trigger for refresh.
type minerUpdateMsg struct{}
type errMsg error

type model struct {
	miners        []*miner_repo.Miner
	selectedMiner *miner_repo.Miner
	updateChan    chan []*miner_repo.Miner
	cursor        int
	err           error

	inProgressFlag chan struct{}
	uiUpdateChan   chan tea.Msg

	mtx         sync.Mutex
	fleetConfig *config.FleetConfig
}

func initialModel(configFilename string) *model {
	updateChan := make(chan []*miner_repo.Miner)
	inProgressFlag := make(chan struct{}, 1)

	fleetConfig, err := config.LoadFleetConfig(configFilename)
	if err != nil {
		log.Fatalf("Failed to load fleet configuration: %v", err)
	}

	m := &model{
		miners:         make([]*miner_repo.Miner, 0),
		updateChan:     updateChan,
		inProgressFlag: inProgressFlag,
		fleetConfig:    fleetConfig,

		uiUpdateChan: make(chan tea.Msg),
	}

	go m.startDataPipeline()

	return m
}

func (m *model) Init() tea.Cmd {
	return m.fetchMiners
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			close(m.updateChan)
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.miners)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.miners) > 0 {
				m.selectedMiner = m.miners[m.cursor]
			}
		case "r":
			return m, m.fetchMiners
		}

	case minerListMsg:
		m.miners = msg
		m.err = nil

	case minerUpdateMsg:
		// This is where we handle the update trigger
		return m, m.fetchMiners

	case errMsg:
		m.err = msg
	}
	return m, nil
}

func (m *model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	s := fmt.Sprintf("Total Miners: %d\n\n", len(m.miners))
	s += "MAC Address | IP Address | Hash Rate | Status\n"
	s += "------------------------------------------------\n"

	// Display up to 10 miners for brevity
	displayCount := len(m.miners)
	if displayCount > 100 {
		displayCount = 100
	}

	for i := 0; i < displayCount; i++ {
		miner := m.miners[i]
		s += fmt.Sprintf("%s | %s | %.2f | %d\n",
			miner.Miner.MacAddress,
			miner.Miner.IPAddress,
			miner.Stats.HashRate,
			miner.Status)
	}

	if len(m.miners) > 100 {
		s += fmt.Sprintf("\n... and %d more miners\n", len(m.miners)-10)
	}

	s += "\nPress 'q' to quit, 'r' to refresh manually\n"
	s += fmt.Sprintf("\nLast error (if any): %v\n", m.err)

	return s
}

func (m *model) fetchMiners() tea.Msg {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return minerListMsg(m.miners)
}

func (m *model) startDataPipeline() {
	ticker := time.NewTicker(time.Duration(INTERVAL_MINS) * time.Minute)
	defer ticker.Stop()

	workerErrors := make(chan error, 1000) // Adjust buffer size as needed

	for {
		select {
		case <-ticker.C:
			fmt.Println("Running scheduled tasks...")
			select {
			case m.inProgressFlag <- struct{}{}:
				go func() {
					defer func() { <-m.inProgressFlag }() // remove the flag
					m.processFleets(workerErrors)
				}()
			default:
				fmt.Println("Skipping the current iteration...")
			}

		case updatedMiners := <-m.updateChan:
			m.mtx.Lock()
			m.miners = updatedMiners
			m.mtx.Unlock()
			m.uiUpdateChan <- minerUpdateMsg{}

		case err := <-workerErrors:
			if err != nil {
				fmt.Println("Error:", err)
				m.mtx.Lock()
				m.err = err
				m.mtx.Unlock()

				m.uiUpdateChan <- minerUpdateMsg{}
			}
		}
	}
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		logFile, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("Failed to set up logging:", err)
			os.Exit(1)
		}
		defer logFile.Close()
	}

	configFilename := "fxhnd.json" // Adjust this to your actual config file name

	m := initialModel(configFilename)
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for msg := range m.uiUpdateChan {
			p.Send(msg)
		}
	}()

	p.Run()
}

func (m *model) processFleets(workerErrors chan<- error) {
	semaphore := make(chan struct{}, MAX_WORKERS)
	var fleetWg sync.WaitGroup

	minersArray := make([]*miner_repo.Miner, 0)

	for fleetIndex, fleet := range m.fleetConfig.Fleets {
		fleetWg.Add(1)
		semaphore <- struct{}{}

		time.Sleep(100 * time.Millisecond)

		go func(fleet config.Fleet) {
			defer fleetWg.Done()
			defer func() { <-semaphore }()

			// log.Printf("Processing fleet: %s\n", fleet.Name)
			log.Printf("Fleet start IP: %s, end IP: %s\n", fleet.Scanner.StartIP, fleet.Scanner.EndIP)

			startIP := net.ParseIP(fleet.Scanner.StartIP)
			endIP := net.ParseIP(fleet.Scanner.EndIP)

			if startIP == nil || endIP == nil {
				workerErrors <- fmt.Errorf("invalid IP address format for fleet %s", fleet.Name)
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

			antMinerCGIModel := make(chan *miner_repo.Miner, len(ips))

			var wg sync.WaitGroup

			for _, ip := range ips {
				time.Sleep(100 * time.Millisecond)
				wg.Add(1)

				clientConnection := http_auth.NewTransport(fleet.Scanner.Username, fleet.Scanner.Password)

				go func(ip net.IP) {
					defer wg.Done()

					newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_system_info.cgi", ip), nil)
					if err != nil {
						return
					}
					resp, err := clientConnection.RoundTrip(newRequest)
					if err != nil {
						return
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return
					}

					var rawGetSystemInfoResponse ant_miner_cgi_queries.RawGetSystemInfoResponse
					if err := json.Unmarshal(body, &rawGetSystemInfoResponse); err != nil {
						return
					}

					antMinerCGIService := ant_miner_cgi_service.NewAntminerCGI(
						&clientConnection,
						miner_domain.Config{
							Username: fleet.Scanner.Username,
							Password: fleet.Scanner.Password,
							Firmware: rawGetSystemInfoResponse.FirmwareType,
						},
						miner_domain.Miner{
							IPAddress:  rawGetSystemInfoResponse.IPAddress,
							MacAddress: rawGetSystemInfoResponse.MACAddress,
						},
						rawGetSystemInfoResponse.MinerType,
					)

					err = antMinerCGIService.CheckStats()
					if err != nil {
						workerErrors <- err
						return
					}

					err = antMinerCGIService.CheckPools()
					if err != nil {
						workerErrors <- err
						return
					}

					newMinerModel := &miner_repo.Miner{
						Miner: miner_domain.Miner{
							IPAddress:  antMinerCGIService.Miner.IPAddress,
							MacAddress: antMinerCGIService.Miner.MacAddress,
						},
						Stats: miner_domain.Stats{
							HashRate:  antMinerCGIService.Stats.HashRate,
							RateIdeal: antMinerCGIService.Stats.RateIdeal,
							Uptime:    antMinerCGIService.Stats.Uptime,
						},
						Config: miner_domain.Config{
							Username: antMinerCGIService.Config.Username,
							Password: antMinerCGIService.Config.Password,
							Firmware: antMinerCGIService.Config.Firmware,
						},
						ModelName: antMinerCGIService.Model,
						Mode:      antMinerCGIService.Mode,

						Status:  antMinerCGIService.Status,
						FleetID: uint(fleetIndex),
					}

					newMinerModel.Fan = make([]int, len(antMinerCGIService.Fan))
					for i, fan := range antMinerCGIService.Fan {
						newMinerModel.Fan[i] = fan.Speed
					}

					newMinerModel.Temperature = make([]int, len(antMinerCGIService.Temperature))
					for i, temp := range antMinerCGIService.Temperature {
						max := 0

						for _, pcbSensor := range temp.PcbSensors {
							if pcbSensor > max {
								max = pcbSensor
							}
						}
						newMinerModel.Temperature[i] = max
					}

					if len(antMinerCGIService.Pools) > 0 {
						newMinerModel.Pools = make([]miner_repo.Pool, len(antMinerCGIService.Pools))
						for i, pool := range antMinerCGIService.Pools {
							newMinerModel.Pools[i] = miner_repo.Pool{
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
						}
					}

					antMinerCGIModel <- newMinerModel

				}(ip)
			}

			wg.Wait()
			close(antMinerCGIModel)

			for minerModel := range antMinerCGIModel {
				miner := &miner_repo.Miner{
					Miner: miner_domain.Miner{
						MacAddress: minerModel.Miner.MacAddress,
						IPAddress:  minerModel.Miner.IPAddress,
					},
				}

				minersArray = append(minersArray, miner)
			}

			fmt.Printf("========================END OF FLEET=========================: %s\n", fleet.Name)
		}(fleet)
	}

	fleetWg.Wait()

	m.updateChan <- minersArray
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
