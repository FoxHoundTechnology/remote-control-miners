import axios from 'axios'
import dotenv from 'dotenv'

// TODO: enum folder
// TODO: store management with zustand

// Load environment variables from .env file
dotenv.config()

const remoteControlServerUrl = process.env.REMOTE_CONTROL_SERVER_URL

type SetMinerModeRequest = {
  MacAddressess: string[]
  Command: Command
}

// NOTE: identical to the domain
export enum Command {
  Normal = 0,
  Sleep = 1,
  LowPower = 2,
  Reboot = 3
}

export const remoteControlAPIService = axios.create({
  baseURL: `${remoteControlServerUrl}/api/miners/control`
})

// TODO: logic to aggregate the macaddress
export const sendCommand = async (setMinerModeRequest: SetMinerModeRequest) => {
  const res = await remoteControlAPIService.post(``, {
    mac_addresses: setMinerModeRequest.MacAddressess,
    command: setMinerModeRequest.Command
  })

  return res.data
}
