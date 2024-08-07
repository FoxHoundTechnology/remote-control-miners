import axios from 'axios'


const remoteControlServerUrl = process.env.REMOTE_CONTROL_SERVER_URL


export const minerAPIService = axios.create({
    baseURL: `${remoteControlServerUrl}/api/miners`,
    headers : {
    'Content-Type': 'application/json'
    }
})

export const fetchMinerLog = async (macAddress: string) => {
    const response = await minerAPIService.post(`/log`, {
        mac_address: macAddress
    })

    return response?.data?.data
}