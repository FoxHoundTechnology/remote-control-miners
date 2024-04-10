package repositories

import (
	"gorm.io/gorm"

	alert "foxhound/internal/application/alert/domain"

	fleet "foxhound/internal/infrastructure/database/repositories/fleet"
)

// TODO: fix schema
// TODO: DTO
type Alert struct {
	gorm.Model
	Alert     alert.Alert              `gorm:"embedded;"`
	Threshold alert.AlertThresholdType `gorm:"type:alert_threshold"`
	Condition alert.AlertConditionType `gorm:"type:alert_condition"`
	Action    alert.AlertActionType    `gorm:"type:alert_action"`
	Layer     alert.AlertLayerType     `gorm:"type:alert_layer"`
	Fleets    []*fleet.Fleet           `gorm:"many2many:alert_fleet;"` // self-referential many-to-many
}
