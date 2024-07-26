package routers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	http_auth "github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	scanner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/scanner/domain"

	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"

	ant_miner_cgi_service "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/service"
)

// TODO: add validation for request object
// TODO: set up a generic response model
// TODO: seggregate the middleware/data conversion/go routines logic into controller
// TODO: variadic function with a map object that comes with different vendor
// TODO: set up an aggregated error response for miner controller logic
// TODO: separate the response logic into controller layer

type MinerRequest struct {
	MacAddress string `json:"mac_address"`
}

type FleetRequest struct {
	FleetID uint `json:"fleet_id"`
}

type MinerControlRequest struct {
	MacAddresses []string             `json:"mac_addresses"`
	Command      miner_domain.Command `json:"command"`
}

func RegisterMinerRoutes(db *gorm.DB, router *gin.Engine) {
	router.POST("/api/miners/log", func(ctx *gin.Context) {
		var minerRequest MinerRequest
		if err := ctx.ShouldBindJSON(&minerRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}

		minerRepository := miner_repo.NewMinerRepository(db)
		miner, err := minerRepository.GetByMacAddress(minerRequest.MacAddress)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error fetching miner data from database",
				"data":    err,
			})
			return
		}

		clientConnection := http_auth.NewTransport(miner.Config.Username, miner.Config.Password)
		antMinerCGIService := ant_miner_cgi_service.NewAntminerCGI(
			&clientConnection,
			miner_domain.Config{
				Username: miner.Config.Username,
				Password: miner.Config.Password,
				Firmware: miner.Config.Firmware,
			},
			miner_domain.Miner{
				IPAddress:  miner.Miner.IPAddress,
				MacAddress: miner.Miner.MacAddress,
			},
			miner.ModelName,
		)

		err = antMinerCGIService.CheckLog()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error fetching miner log",
				"data":    err,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully fetched miner log",
			"data":    antMinerCGIService.Log,
		})

	})

	router.POST("/api/miners/detail", func(ctx *gin.Context) {
		var minerRequest MinerRequest
		if err := ctx.ShouldBindJSON(&minerRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}

		minerRepository := miner_repo.NewMinerRepository(db)
		miner, err := minerRepository.GetByMacAddress(minerRequest.MacAddress)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error fetching miner detail",
				"data":    err,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully fetched miner detail",
			"data":    miner,
		})
	})

	router.GET("/api/miners/list", func(ctx *gin.Context) {
		minerRepository := miner_repo.NewMinerRepository(db)
		miners, err := minerRepository.List()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error fetching miners",
				"data":    err,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully fetched miners",
			"data":    miners,
		})
	})

	router.GET("/api/miners/fleets", func(ctx *gin.Context) {

		var fleetRequest FleetRequest
		if err := ctx.ShouldBindJSON(&fleetRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}

		minerRepository := miner_repo.NewMinerRepository(db)
		miners, err := minerRepository.ListByFleetID(fleetRequest.FleetID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error fetching miners",
				"data":    err,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully fetched miners",
			"data":    miners,
		})
	})

	// TODO: fleet_id
	// TODO: seggregate the caller logic from the router endpoint to the controller folder
	router.POST("/api/miners/control", func(ctx *gin.Context) {

		var minerControlRequest MinerControlRequest
		if err := ctx.ShouldBindJSON(&minerControlRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}

		minerRepository := miner_repo.NewMinerRepository(db)
		miners, err := minerRepository.ListByMacAddresses(minerControlRequest.MacAddresses)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error sending config requests to miners",
				"data":    err,
			})
		}

		antMinerCGIServiceArray := []ant_miner_cgi_service.AntminerCGI{}

		// FIXME: data access to username and password
		// clientConnection := http_auth.NewTransport(fleet.Scanner.Config.Username, fleet.Scanner.Config.Password)

		for _, miner := range miners {
			minerType := miner.MinerType

			// TODO: fix the logic to retrieve username/password
			clientConnection := http_auth.NewTransport(miner.Config.Username, miner.Config.Password)

			switch minerType {
			case scanner_domain.AntminerCgi:
				antMinerCGIService := ant_miner_cgi_service.NewAntminerCGI(
					&clientConnection,
					miner_domain.Config{
						Username: miner.Config.Username,
						Password: miner.Config.Password,
						Firmware: miner.Config.Firmware,
					},
					miner_domain.Miner{
						IPAddress:  miner.Miner.IPAddress,
						MacAddress: miner.Miner.MacAddress,
					},
					miner.ModelName,
				)
				antMinerCGIServiceArray = append(antMinerCGIServiceArray, *antMinerCGIService)

			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "miner type not supported",
				})
			}
		}

		minerControlerErrorChannel := make(chan error)
		var wg sync.WaitGroup

		for _, antMinerCGIService := range antMinerCGIServiceArray {
			wg.Add(1)
			go func(antMinerCGIService ant_miner_cgi_service.AntminerCGI) {
				defer wg.Done()
				switch minerControlRequest.Command {

				case miner_domain.Normal:
					err := antMinerCGIService.SetNormalMode()
					if err != nil {
						minerControlerErrorChannel <- err
					}

				case miner_domain.Sleep:
					err := antMinerCGIService.SetSleepMode()
					if err != nil {
						minerControlerErrorChannel <- err
					}

				case miner_domain.LowPower:
					err := antMinerCGIService.SetLowPowerMode()
					if err != nil {
						minerControlerErrorChannel <- err
					}

				case miner_domain.Reboot:
					err := antMinerCGIService.Reboot()
					if err != nil {
						minerControlerErrorChannel <- err
					}

				default:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": "mode not supported",
					})
				}
			}(antMinerCGIService)
		}

		// NOTE: here goes the aggregated error response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully updated the miners",
		})
	})
}
