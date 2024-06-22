// NOTE: the response object from the server is
//       in snake case
export type Scanner = {
  fleetName: string
  client: string | null

  ipRange: string | null

  startIp: string | null
  endIp: string | null

  minerType: string | null
  mins: number | null
  password: string | null
  username: string | null

  isActive?: boolean
  lastRun?: Date
  lastSuccess?: Date
  lastFailure?: Date
  errorMessage?: string

  alertBelowHashRate?: number | null
  alertHashrateEnabled?: boolean | null
  warnAboveHashRate?: number | null
  warnHashrateEnabled?: boolean | null

  alertAboveFan?: number | null
  alertFanEnabled?: boolean | null
  warnAboveFan?: number | null
  warnFanEnabled?: boolean | null

  alertAboveTemp?: number | null
  alertTempEnabled?: boolean | null
  warnAboveTemp?: number | null
  warnTempEnabled?: boolean | null
}

