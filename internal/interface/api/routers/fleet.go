package routers

import (
	"net/http"

	fleet_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/fleet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterFleetRoutes(db *gorm.DB, router *gin.Engine) {

	router.GET("/fleets/list", func(ctx *gin.Context) {
		// scanners := []scanner.Scanner{}
		fleets := []fleet_repo.Fleet{}
		db.Find(&fleets)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "fleet list test",
			"data":    fleets,
		})
	})
}
