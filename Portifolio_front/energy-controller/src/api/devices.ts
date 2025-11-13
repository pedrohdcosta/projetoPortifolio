import api from './axios'

export interface Device {
  id: number
  name: string
  room: string
}

/**
 * Fetches the list of devices from the API
 * @returns Promise<Device[]>
 */
export async function listDevices(): Promise<Device[]> {
  const response = await api.get<Device[]>('/devices')
  return response.data
}
