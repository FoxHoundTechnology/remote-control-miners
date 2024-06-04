package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	scanner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/scanner"
)

// TODO: separate the response logic into controller layer
func RegisterScannerRoutes(db *gorm.DB, router *gin.Engine) {

	router.GET("/api/scanners/list", func(ctx *gin.Context) {
		scanners := []scanner_repo.Scanner{}
		db.Find(&scanners)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "scanner list test",
			"data":    scanners,
		})
	})

}
