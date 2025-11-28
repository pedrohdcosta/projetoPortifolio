const STORAGE_KEY = 'device_simulation_mode_v1'

function loadAll(): Record<string, boolean> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return {}
    return JSON.parse(raw)
  } catch (e) {
    console.warn('Failed to load simulation modes', e)
    return {}
  }
}

function saveAll(map: Record<string, boolean>) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(map))
  } catch (e) {
    console.warn('Failed to save simulation modes', e)
  }
}

export function getSimulationMode(deviceId: number): boolean | undefined {
  const map = loadAll()
  return map[String(deviceId)]
}

export function setSimulationMode(deviceId: number, enabled: boolean) {
  const map = loadAll()
  map[String(deviceId)] = !!enabled
  saveAll(map)
}

export function toggleSimulationMode(deviceId: number) {
  const map = loadAll()
  const key = String(deviceId)
  map[key] = !map[key]
  saveAll(map)
  return map[key]
}

export function getAllSimulationModes(): Record<string, boolean> {
  return loadAll()
}

export default { getSimulationMode, setSimulationMode, toggleSimulationMode, getAllSimulationModes }
