import axios from 'axios'
import dotenv from 'dotenv'

import { MinerType } from 'src/types/apps/minerTypes'

// TODO: enum folder
// TODO: request/response folder (redo the one above)
// TODO: store management with zustand

// TODO: -- DUPLICATES --
const convertResponse = (response: any): MinerType[] => {
  try {
    return response.map((item: any) => ({
      id: item.ID,
      macAddress: item.Miner?.MacAddress || 'N/A',
      ip: item.Miner?.IPAddress || 'N/A',
      minerType: item.MinerType || 'N/A',
      serialNumber: 'N/A', // Example placeholder, adjust as per actual data
      status:
        item.Status === 0
          ? 'Online'
          : item.Status === 1
          ? 'Offline'
          : item.Status === 2
          ? 'Disabled'
          : item.Status === 3
          ? 'Hashrate Error' // ctd...
          : item.Status === 4
          ? 'Warning'
          : 'N/A',
      mode: item.Mode === 0 ? 'Normal' : item.Mode === 1 ? 'Sleep' : item.Mode === 2 ? 'Low Power' : 'N/A',
      model: item?.ModelName || 'N/A',
      hashRate: item.Stats?.HashRate || 0,
      rateIdeal: item.Stats?.RateIdeal || 0,
      maxFan: item.Fan?.length > 0 ? Math.max(...item.Fan) : 0,
      maxTemp: item.Temperature?.length > 0 ? Math.max(...item.Temperature) : 0,
      fanArr: item.Fan || [],
      tempArr: item.Temperature || [],
      upTime: item.Stats?.Uptime || 0,
      location: 'N/A', // Placeholder, add actual logic if location data exists
      firmware: item.Config?.Firmware || 'N/A',
      client: 'N/A', // Placeholder
      fleetName: 'N/A', // Placeholder
      lastUpdated: item.UpdatedAt,
      ipRange: 'N/A', // Placeholder
      password: item.Config?.Password || 'N/A',
      username: item.Config?.Username || 'N/A',
      // alert/warning config
      alertAboveHashRate: item.alert_above_hash_rate || null,
      alertBelowHashRate: item.alert_below_hash_rate || null,
      alertHashrateEnabled: item.alert_hashrate_enabled || null,
      warnAboveHashRate: item.warn_above_hash_rate || null,
      warnBelowHashRate: item.warn_below_hash_rate || null,
      warnHashrateEnabled: item.warn_hashrate_enabled || null,
      alertBelowFan: item.alert_below_fan || null,
      alertAboveFan: item.alert_above_fan || null,
      alertFanEnabled: item.alert_fan_enabled || null,
      warnAboveFan: item.warn_above_fan || null,
      warnBelowFan: item.warn_below_fan || null,
      warnFanEnabled: item.warn_fan_enabled || null,
      alertBelowTemp: item.alert_below_temp || null,
      alertAboveTemp: item.alert_above_temp || null,
      alertTempEnabled: item.alert_temp_enabled || null,
      warnBelowTemp: item.warn_below_temp || null,
      warnAboveTemp: item.warn_above_temp || null,
      warnTempEnabled: item.warn_temp_enabled || null,
      warnMessage: item.warn_message || 'N/A',
      activityLogs: item.activityLogs || []
    }))
  } catch (err) {
    console.log('error converting response: ', err)
    return []
  }
}
// -- THE END OF DUPLICATES --
dotenv.config()

// Get REMOTE_CONTROL_SERVER_URL from environment variables
const remoteControlServerUrl = process.env.REMOTE_CONTROL_SERVER_URL

export const remoteControlAPIService = axios.create({
  baseURL: `${remoteControlServerUrl}/api/miners`,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const fetchMinerInfo = async (macAddress: string) => {
  const response = await remoteControlAPIService.post(`/detail`, { mac_address: macAddress })
  // NOTE: insert only one
  const convertedResponse = convertResponse([response?.data?.data])
  return convertedResponse[0]
}

export const fetchMinerStats = async (
  macAddress: string,
  interval: number,
  intervalUnit: string,
  window: number,
  windowUnit: string
) => {
  const response = await remoteControlAPIService.post(`/timeseries/minerstats`, {
    mac_address: macAddress,
    interval: interval,
    interval_unit: intervalUnit,
    window: window,
    window_unit: windowUnit
  })

  return response.data?.data
}

export const fetchPoolStats = async (
  macAddress: string,
  interval: number,
  intervalUnit: string,
  window: number,
  windowUnit: string
) => {
  const response = await remoteControlAPIService.post(`/timeseries/poolstats`, {
    mac_address: macAddress,
    interval: interval,
    interval_unit: intervalUnit,
    window: window,
    window_unit: windowUnit
  })

  return response.data?.data
}
