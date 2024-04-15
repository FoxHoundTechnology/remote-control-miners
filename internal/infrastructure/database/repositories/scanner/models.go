package repositories

import (
	"gorm.io/gorm"

	domain "foxhound/internal/application/scanner/domain"
)

type Scanner struct {
	gorm.Model
	Scanner domain.Scanner `gorm:"embedded;"`
	Config  domain.Config  `gorm:"embedded;"`

	MinerType string `gorm:"type:miner_type"`
	Owner     string
	Alerts    []*Alert `gorm:"many2many:scanner_alerts;"`
}

type Alert struct {
	gorm.Model
	Name      string
	Threshold domain.AlertThresholdType `gorm:"comment:'ThresholdCount=0, ThresholdRate=1'"`
	Condition domain.AlertConditionType `gorm:"comment:'Hashrate=0, Temperature=1, FanSpeed=2, PoolShares=3, OfflineMiners=4, MissingHashboards=5'"`
	Action    domain.AlertActionType    `gorm:"comment:'Reboot=0, Sleep=1, Normal=2, ChangePool=3'"`
	Layer     domain.AlertLayerType     `gorm:"comment:'InfoAlert=0, WarningAlert=1, ErrorAlert=2, FataltAlert=3'"`
	State     domain.AlertState         `gorm:"comment:'Received=0, InProgress=1, Completed=2'"`
	Log       []AlertLog                `gorm:"foreignKey:AlertID;references:ID"`
	// (CreatedAt): a field in gorm.Model.
	Scanners []*Scanner `gorm:"many2many:scanner_alerts;"`
}

type AlertLog struct {
	gorm.Model
	Log     domain.Log `gorm:"embedded;"`
	AlertID uint       `gorm:"foreignKey:AlertID;"`
}
