package routers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	scanner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/scanner/domain"

	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"

	ant_miner_cgi_service "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/service"
)

// TODO: add validation for request object
// TODO: set up universal response model
// TODO: seggregate the middleware/data conversion/go routines logic into controller
// TODO: variadic function with a map object that comes with different vendor
// TODO: set up an aggregated error response for miner controller logic
// TODO: separate the response logic into controller layer

type MinerControlRequest struct {
	MacAddresses []string          `json:"mac_addresses"`
	Mode         miner_domain.Mode `json:"mode"`
}

func RegisterMinerRoutes(db *gorm.DB, router *gin.Engine) {

	router.GET("/miners/list", func(ctx *gin.Context) {
		minerRepository := miner_repo.NewMinerRepository(db)
		miners, err := minerRepository.List()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error fetching miners",
				"data":    err,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully fetched miners",
			"data":    miners,
		})
	})

	router.POST("/miners/control", func(ctx *gin.Context) {

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

		for _, miner := range miners {
			minerType := miner.MinerType

			switch minerType {
			case scanner_domain.AntminerCgi:
				antMinerCGIService := ant_miner_cgi_service.NewAntminerCGI(
					miner_domain.Config{
						Username: miner.Config.Username,
						Password: miner.Config.Password,
						Firmware: miner.Config.Firmware,
					},
					miner_domain.Miner{
						IPAddress:  miner.Miner.IPAddress,
						MacAddress: miner.Miner.MacAddress,
					},
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
				switch minerControlRequest.Mode {

				case miner_domain.NormalMode:
					err := antMinerCGIService.SetNormalMode()
					if err != nil {
						minerControlerErrorChannel <- err
					}

				case miner_domain.SleepMode:
					err := antMinerCGIService.SetSleepMode()
					if err != nil {
						minerControlerErrorChannel <- err
					}

				case miner_domain.LowPowerMode:
					err := antMinerCGIService.SetLowPowerMode()
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
