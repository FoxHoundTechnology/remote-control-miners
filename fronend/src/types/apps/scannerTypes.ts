import { ThemeColor } from 'src/@core/layouts/types'

// snake case
export type ScannerType = {
  id?: number | string // = fleet_name
  fleet_name: string
  client: string
  start_ip: string
  end_ip: string
  miner_type: string
  mins: number
  password: string
  username: string
  is_active: boolean
  last_run: string
  last_success: string
  last_failure: string
  error_message: string
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
}
