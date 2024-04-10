package alert

type AlertThresholdType string
type AlertConditionType string
type AlertActionType string
type AlertLayerType uint8

const (
	ThresholdCount AlertThresholdType = "count"
	ThresholdRate  AlertThresholdType = "rate"
)

const (
	Hashrate          AlertConditionType = "hashrate"
	Temperature       AlertConditionType = "temperature"
	FanSpeed          AlertConditionType = "fan_speed"
	PoolShares        AlertConditionType = "pool_shares"
	OfflineMiners     AlertConditionType = "offline_miners"
	MissingHashboards AlertConditionType = "missing_hashboards"
)

const (
	Reboot     AlertActionType = "reboot"
	Sleep      AlertActionType = "sleep"
	Normal     AlertActionType = "normal" // turn on
	ChangePool AlertActionType = "change_pool"
)

const (
	InfoAlert    AlertLayerType = 1
	WarningAlert AlertLayerType = 2
	ErrorAlert   AlertLayerType = 3
	FataltAlert  AlertLayerType = 4
)
