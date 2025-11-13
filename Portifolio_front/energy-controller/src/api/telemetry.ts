import api from './axios'

export interface TelemetryRow {
  time: string
  power_w: number
  voltage: number
  current: number
}

export interface TelemetryParams {
  device_id?: number
  from?: string
  to?: string
  limit?: number
}

/**
 * Fetches telemetry data from the API
 * @param params Query parameters for the telemetry request
 * @returns Promise<TelemetryRow[]>
 */
export async function fetchTelemetry(params: TelemetryParams): Promise<TelemetryRow[]> {
  const response = await api.get<TelemetryRow[]>('/telemetry', { params })
  return response.data
}
