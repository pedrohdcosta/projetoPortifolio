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
