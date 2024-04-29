package repositories

import (
	"foxhound/internal/application/scanner/domain"

	"gorm.io/gorm"
)

// TODO: scanner name should be unique
type Scanner struct {
	gorm.Model
	Name      string
	Scanner   domain.Scanner   `gorm:"embedded;"`
	Config    domain.Config    `gorm:"embedded;"`
	MinerType domain.MinerType `gorm:"comment:'AntMinerCgi=0'"`
	Owner     string

	Alert Alert

	FleetID uint
}

type Alert struct {
	gorm.Model
	Name   string
	Action domain.AlertActionType `gorm:"comment:'Reboot=0, Sleep=1, Normal=2, ChangePool=3'"`
	State  domain.AlertState      `gorm:"comment:'Monitoring=0, Triggered=1, Resolving=2, Resolved=3'"`

	Condition []AlertCondition `gorm:"foreignKey:AlertID;references:ID"`
	Log       []AlertLog       `gorm:"foreignKey:AlertID;references:ID"`

	ScannerID uint
}

type AlertCondition struct {
	gorm.Model
	Value     domain.AlertValue
	Threshold domain.AlertThresholdType `gorm:"comment:'ThresholdCount=0, ThresholdRate=1'"`
	Condition domain.AlertConditionType `gorm:"comment:'Hashrate=0, Temperature=1, FanSpeed=2, PoolShares=3, OfflineMiners=4, MissingHashboards=5'"`
	Layer     domain.AlertLayerType     `gorm:"comment:'InfoAlert=0, WarningAlert=1, ErrorAlert=2, FataltAlert=3'"`

	AlertID uint
}

type AlertLog struct {
	gorm.Model
	Log domain.Log

	AlertID uint
}
