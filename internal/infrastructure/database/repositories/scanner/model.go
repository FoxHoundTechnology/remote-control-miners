package repositories

import (
	"github.com/FoxHoundTechnology/remote-control-miners/internal/application/scanner/domain"

	"gorm.io/gorm"
)

// TODO: add the logic to prevent the same conditions type
// TODO: test cascade condition for onUpdate/onDelete
//	(e.g. Hashrate and Hashrate) from being added to the same alert

type Scanner struct {
	gorm.Model
	Name      string
	Scanner   domain.Scanner   `gorm:"embedded;"`
	Config    domain.Config    `gorm:"embedded;"`
	MinerType domain.MinerType `gorm:"comment:'AntMinerCgi=0'"`
	Owner     string

	Alert Alert `gorm:"onDelete:CASCADE; onUpdate:CASCADE"`

	FleetID uint
}

type Alert struct {
	gorm.Model
	Name   string
	Action domain.AlertActionType `gorm:"comment:'Reboot=0, Sleep=1, Normal=2'"`
	State  domain.AlertState      `gorm:"comment:'Monitoring=0, Triggered=1, Resolving=2, Resolved=3'"`
	Active bool                   `gorm:"default:true"`

	Condition []AlertCondition `gorm:"foreignKey:AlertID;references:ID"`
	Log       []AlertLog       `gorm:"foreignKey:AlertID;references:ID"`

	ScannerID uint
}

type AlertCondition struct {
	gorm.Model
	TriggerValue  domain.AlertTriggerValue  `gorm:"comment:'Value for threshold (e.g. 100 TH/s, 80°C, and 100 RPM)'"`
	MachineCount  domain.AlertMachineCount  `gorm:"comment:'Number of machines for alert to get triggerred'"`
	ThresholdType domain.AlertThresholdType `gorm:"comment:'ThresholdCount=1, ThresholdRate=1'"`
	ConditionType domain.AlertConditionType `gorm:"comment:'Hashrate=0, Temperature=1, FanSpeed=2, PoolShares=3, OfflineMiners=4, MissingHashboards=5'"`
	LayerType     domain.AlertLayerType     `gorm:"comment:'InfoAlert=0, WarningAlert=1, ErrorAlert=2, FataltAlert=3'"`

	AlertID uint
}

type AlertLog struct {
	gorm.Model
	Log domain.Log

	AlertID uint
}
