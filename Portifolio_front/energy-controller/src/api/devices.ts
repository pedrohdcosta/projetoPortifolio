import api from "./axios";

export interface Device {
  id: number;
  user_id: number;
  name: string;
  room?: string;
  type?: string;
  status?: "online" | "offline";
  metadata?: string;
  created_at?: string;
  last_seen?: string;
}

export interface TelemetryData {
  id: number;
  device_id: number;
  timestamp: string;
  power: number;
  voltage?: number;
  current?: number;
}

export interface TelemetrySummary {
  device_id: number;
  period: string;
  start_time: string;
  end_time: string;
  total_records: number;
  avg_power: number;
  max_power: number;
  min_power: number;
  total_energy: number;
  avg_voltage?: number;
  avg_current?: number;
}

export interface CreateDeviceRequest {
  name: string;
  room?: string;
  type?: string;
  metadata?: string;
}

export interface UpdateDeviceRequest {
  name?: string;
  room?: string;
  type?: string;
  status?: string;
  metadata?: string;
}

export interface CreateTelemetryRequest {
  device_id: number;
  power: number;
  voltage?: number;
  current?: number;
}

// Device API functions
export async function listDevices(): Promise<Device[]> {
  const { data } = await api.get("/devices");
  return data;
}

export async function getDevice(id: number): Promise<Device> {
  const { data } = await api.get(`/devices/${id}`);
  return data;
}

export async function createDevice(req: CreateDeviceRequest): Promise<Device> {
  const { data } = await api.post("/devices", req);
  return data;
}

export async function updateDevice(id: number, req: UpdateDeviceRequest): Promise<Device> {
  const { data } = await api.put(`/devices/${id}`, req);
  return data;
}

export async function deleteDevice(id: number): Promise<void> {
  await api.delete(`/devices/${id}`);
}

// Telemetry API functions
export async function listTelemetry(deviceId?: number, limit = 100): Promise<TelemetryData[]> {
  const params: Record<string, any> = { limit };
  if (deviceId) {
    params.device_id = deviceId;
  }
  const { data } = await api.get("/telemetry", { params });
  return data;
}

export async function fetchTelemetry(deviceId: number, limit = 100): Promise<TelemetryData[]> {
  const { data } = await api.get("/telemetry", {
    params: { device_id: deviceId, limit },
  });
  return data;
}

export async function createTelemetry(req: CreateTelemetryRequest): Promise<TelemetryData> {
  const { data } = await api.post("/telemetry", req);
  return data;
}

export async function deleteTelemetry(id: number): Promise<void> {
  await api.delete(`/telemetry/${id}`);
}

// New API functions for device-telemetry integration

// Get latest telemetry for all user's devices (one reading per device)
export async function getLatestTelemetry(): Promise<TelemetryData[]> {
  const { data } = await api.get("/telemetry/latest");
  return data;
}

// Get telemetry for a specific device using the dedicated endpoint
export async function getDeviceTelemetry(deviceId: number, limit = 100): Promise<TelemetryData[]> {
  const { data } = await api.get(`/devices/${deviceId}/telemetry`, { params: { limit } });
  return data;
}

// Get telemetry summary for a device (aggregated data)
export async function getDeviceTelemetrySummary(
  deviceId: number, 
  period: 'day' | 'week' | 'month' = 'day'
): Promise<TelemetrySummary> {
  const { data } = await api.get(`/devices/${deviceId}/telemetry/summary`, { params: { period } });
  return data;
}

// Get the latest telemetry reading for a specific device
export async function getDeviceLatestTelemetry(deviceId: number): Promise<TelemetryData> {
  const { data } = await api.get(`/devices/${deviceId}/telemetry/latest`);
  return data;
}
