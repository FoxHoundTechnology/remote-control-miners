package repositories

import (
	"gorm.io/gorm"

	domain "foxhound/internal/application/alert/domain"
)

type Alert struct {
	gorm.Model
	Alert     domain.Alert              `gorm:"embedded;"`
	Threshold domain.AlertThresholdType `gorm:"comment:'ThresholdCount=0, ThresholdRate=1'"`
	Condition domain.AlertConditionType `gorm:"comment:'Hashrate=0, Temperature=1, FanSpeed=2, PoolShares=3, OfflineMiners=4, MissingHashboards=5'"`
	Action    domain.AlertActionType    `gorm:"comment:'Reboot=0, Sleep=1, Normal=2, ChangePool=3'"`
	Layer     domain.AlertLayerType     `gorm:"comment:'InfoAlert=0, WarningAlert=1, ErrorAlert=2, FataltAlert=3'"`
}
