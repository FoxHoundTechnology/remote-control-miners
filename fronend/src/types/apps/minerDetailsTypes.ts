// TODO: FIX ME
export type MinerDetailsData = {
  name: string
  ip: string
  mac_address: string
  miner_type: string
  serial_number: string
  status: string
  mode: string
  hash_rate: number
  rate_ideal: number
  max_fan: number[]
  max_temp: number[]
  up_time: number
  location: string
  firmware: string
  client: string
  fleet_name: string
  last_updated: string
  ip_range: string
  password: string
  username: string
  alert_below_hash_rate: number
  alert_hashrate_enabled: boolean
  warn_below_hash_rate: number
  warn_hashrate_enabled: boolean
  alert_above_fan: number
  alert_fan_enabled: boolean
  warn_above_fan: number
  warn_fan_enabled: boolean
  alert_above_temp: number
  alert_temp_enabled: boolean
  warn_above_temp: number
  warn_temp_enabled: boolean
  warning_message: string
  // activity_logs : string[]
}

// For miner details page
// in snake case
export type MinerTimeSeriesData = {
  hashrate: number
  temp_sensor: number[] // Assuming a maximum of 10 temperature sensors
  fan_sensor: number[] // Assuming a maximum of 10 fan sensors
}

export type MinerTimeSeriesDataResponse = {
  miner_time_series_record: MinerTimeSeriesData[]
  timestamps: Date[] // Using Date for time.Time
}

export type PoolStatsData = {
  accepted: number
  rejected: number
  stale: number
}

export type PoolTimeSeriesDataResponse = {
  pool_time_series_record: PoolStatsData[]
  timestamps: Date[] // Using Date for time.Time
}
