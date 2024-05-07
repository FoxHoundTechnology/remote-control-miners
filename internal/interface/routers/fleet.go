package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	fleet_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/fleet"
)

// TODO: separate the response logic into controller layer

func RegisterFleetRoutes(db *gorm.DB, router *gin.Engine) {

	router.GET("/fleets/list", func(ctx *gin.Context) {
		fleetRepository := fleet_repo.NewFleetRepository(db)
		fleets, err := fleetRepository.ListScannersByFleet()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error",
				"data":    err,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "fleet list test",
			"data":    fleets,
		})
	})
}
