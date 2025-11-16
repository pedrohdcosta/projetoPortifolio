import api from "./axios";

export interface Device {
  id: number;
  name: string;
  type: string;
  status: "online" | "offline";
  lastSeen?: string;
}

export interface TelemetryData {
  id: number;
  deviceId: number;
  timestamp: string;
  power: number;
  voltage?: number;
  current?: number;
}

export async function listDevices(): Promise<Device[]> {
  const { data } = await api.get("/devices");
  return data;
}

export async function fetchTelemetry(deviceId: number, limit = 100): Promise<TelemetryData[]> {
  const { data } = await api.get(`/devices/${deviceId}/telemetry`, {
    params: { limit },
  });
  return data;
}
