package routers

import (
	"net/http"

	// NOTE: it includes the time series repository as well
	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"

	timeseries_database "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/influxdb"

	"github.com/gin-gonic/gin"
)

// TODO: add validation for request object
// TODO: set up universal response model
// TODO: seggregate the middleware/data conversion/go routines logic into controller
// TODO: variadic function with a map object that comes with different vendor
// TODO: set up an aggregated error response for miner controller logic
// TODO: separate the response logic into controller layer

type MinerTimeSeriesRequest struct {
	MacAddress string `json:"mac_addresses"`
	Interval   int    `json:"interval"`
}

func RegisterMinerTimeSeriesRoutes(router *gin.Engine) {

	InfluxDBConnectionSettings := timeseries_database.Init()
	minerTimeSeriesRepository := miner_repo.NewMinerTimeSeriesRepository(InfluxDBConnectionSettings)

	router.GET("/miners/timeseries/minerstats", func(ctx *gin.Context) {
		request := MinerTimeSeriesRequest{}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}
		res, err := minerTimeSeriesRepository.ReadMinerData(request.MacAddress, request.Interval)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error getting miner data",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully fetched miner timeseries data",
			"data":    res,
		})
	})

	router.GET("miners/timeseries/poolstats", func(ctx *gin.Context) {
		requests := MinerTimeSeriesRequest{}
		if err := ctx.ShouldBindJSON(&requests); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Incorrect request object": err.Error()})
			return
		}
		res, err := minerTimeSeriesRepository.ReadPoolData(requests.MacAddress, requests.Interval)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error getting pool data",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully fetched pool timeseries data",
			"data":    res,
		})

	})

}
