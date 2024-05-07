package routers

import (
	"net/http"

	miner_repo "github.com/FoxHoundTechnology/remote-control-miners/internal/infrastructure/database/repositories/miner"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMinerRoutes(db *gorm.DB, router *gin.Engine) {

	router.GET("/miners/list", func(ctx *gin.Context) {

		miners := []miner_repo.Miner{}
		db.Find(&miners)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "fleet list test",
			"data":    miners,
		})
	})
}
