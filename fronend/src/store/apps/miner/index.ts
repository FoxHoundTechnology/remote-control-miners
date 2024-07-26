import axios from 'axios'
import { MinerType } from 'src/types/apps/minerTypes'

const remoteControlServerUrl = process.env.REMOTE_CONTROL_SERVER_URL

export const minerAPIService = axios.create({
  baseURL: `${remoteControlServerUrl}/api/miners`
})
// TODO: store management with zustand
// TODO: enum map fo mode/status (with falsy values supports)
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
      username: item.Config?.Username || 'N/A'
    }))
  } catch (err) {
    console.log('error converting response: ', err)
    return []
  }
}
// -- THE END OF DUPLICATES --

export const fetchMinerList = async () => {
  const response = await minerAPIService.get('/list')
  console.log('res.data: ', response)
  const convertedResponse = convertResponse(response?.data?.data) ?? []
  return { miners: convertedResponse }
}

// ========= Util functions for miner data =========
export const ExtractFields = (data: MinerType[]): { [key: string]: string[] } => {
  const fields: (keyof MinerType)[] = ['client', 'minerType', 'status', 'location']
  const result: { [key: string]: string[] } = {}
  data?.forEach(item => {
    fields.forEach(field => {
      // Using a type assertion to tell TypeScript that `field` is a key of `MinerType`
      const key = field as keyof MinerType
      if (!result[field]) result[field] = []
      // Need to make sure that the value in MinerType[field] is a string.
      if (item[key] !== 'N/A' && !result[field].includes(item[key] as string)) {
        result[field].push(item[key] as string)
      }
    })
  })

  return result
}
