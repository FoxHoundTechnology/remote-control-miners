// NOTE: the response object from the server is
//       in snake case
// TODO: will fix this later
export type MinerType = {
  id: number
  macAddress: string
  ip: string
  name: string
  minerType: string
  serialNumber: string
  status: string
  mode: string
  model: string
  hashRate: number
  maxFan: number[]
  maxTemp: number[]
  fanArr?: number[]
  tempArr?: number[]
  upTime: number
  location: string
  firmware: string
  client: string
  fleetName: string
  lastUpdated: string
  ipRange: string
  password: string
  username: string
  // alertAboveHashRate: number | null
  // alertBelowHashRate: number | null
  // alertHashrateEnabled: boolean | null
  // warnAboveHashRate: number | null
  // warnBelowHashRate: number | null
  // warnHashrateEnabled: boolean | null
  // alertBelowFan: number | null
  // alertAboveFan: number | null
  // alertFanEnabled: boolean | null
  // warnAboveFan: number | null
  // warnBelowFan: number | null
  // warnFanEnabled: boolean | null
  // alertBelowTemp: number | null
  // alertAboveTemp: number | null
  // alertTempEnabled: boolean | null
  // warnBelowTemp: number | null
  // warnAboveTemp: number | null
  // warnTempEnabled: boolean | null
  // warnMessage: string | null
  // activityLogs: string[]
}
