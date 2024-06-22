// Enums for TargetType and ActionType
export enum TargetType {
  Hashrate = 'hashrate',
  Temperature = 'temperature',
  FanSpeed = 'fan_speed',
  // Shares = 'shares',
  Offline = 'offline',
  MissingHashboard = 'missing_hashboard'
  // PoolConfig = 'pool_config'
}

export enum ActionType {
  Reboot = 'reboot',
  SleepMode = 'sleep_mode',
  NormalMode = 'normal_mode'
  // ChangePool = 'change_pool'
}

// Trigger Model
export type Trigger = {
  ID: number
  name: string
  user_id?: number
  interval: number
  active: boolean
  last_executed: Date
  targets: Target[]
  actions: Action[]
  histories?: TriggerHistory[] | null
}

// Target Model
export type Target = {
  ID?: number
  type: TargetType
  percentage?: number
  value?: number
  trigger_id?: number
}

// Action Model
export type Action = {
  ID?: number
  type: ActionType
  value: string
  interval: number
  trigger_id?: number
}

export type TriggerHistory = {
  ID: number
  timestamp: Date
  message: string
  trigger_id: number
}
