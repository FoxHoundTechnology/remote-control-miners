package routers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

type MinerDetailRequest struct {
	MacAddress string `json:"mac_address"`
}

type MinerControlRequest struct {
	MacAddresses []string             `json:"mac_addresses"`
	Command      miner_domain.Command `json:"command"`
}

func RegisterMinerRoutes(db *gorm.DB, router *gin.Engine) {

	router.POST("/api/miners/detail", func(ctx *gin.Context) {
		var minerDetailRequest MinerDetailRequest
		if err := ctx.ShouldBindJSON(&minerDetailRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}

		minerRepository := miner_repo.NewMinerRepository(db)
		miner, err := minerRepository.GetByMacAddress(minerDetailRequest.MacAddress)
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
		}

		ctx.Header("Cache-Control", "public, max-age="+strconv.Itoa(5))

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully fetched miners",
			"data":    miners,
		})
	})

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
