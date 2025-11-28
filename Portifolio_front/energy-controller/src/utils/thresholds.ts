export interface DeviceThresholds {
  warning: number
  danger: number
}

const STORAGE_KEY = 'device_thresholds_v1'

function loadAll(): Record<string, DeviceThresholds> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return {}
    return JSON.parse(raw)
  } catch (e) {
    console.warn('Failed to load thresholds', e)
    return {}
  }
}

function saveAll(map: Record<string, DeviceThresholds>) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(map))
  } catch (e) {
    console.warn('Failed to save thresholds', e)
  }
}

export function getThresholds(deviceId: number): DeviceThresholds | undefined {
  const map = loadAll()
  return map[String(deviceId)]
}

export function setThresholds(deviceId: number, thresholds: DeviceThresholds) {
  const map = loadAll()
  map[String(deviceId)] = thresholds
  saveAll(map)
}

export function removeThresholds(deviceId: number) {
  const map = loadAll()
  delete map[String(deviceId)]
  saveAll(map)
}

export function getAllThresholds(): Record<string, DeviceThresholds> {
  return loadAll()
}

export default { getThresholds, setThresholds, removeThresholds, getAllThresholds }
