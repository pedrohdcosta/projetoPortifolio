import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Devices from '../Devices.vue'
import * as devicesApi from '../../api/devices'
import * as simulationMode from '../../utils/simulationMode'

// Mock API calls
vi.mock('../../api/devices', () => ({
  listDevices: vi.fn(),
  getLatestTelemetry: vi.fn(),
  testDeviceConnection: vi.fn(),
  toggleDevice: vi.fn(),
  createDevice: vi.fn(),
  deleteDevice: vi.fn(),
  updateDevice: vi.fn(),
  getDeviceTelemetry: vi.fn(),
  getDeviceTelemetrySummary: vi.fn(),
  simulateTelemetry: vi.fn(),
  simulateBulkTelemetry: vi.fn(),
}))

vi.mock('../../utils/simulationMode', () => ({
  getAllSimulationModes: vi.fn(),
  setSimulationMode: vi.fn(),
  toggleSimulationMode: vi.fn(),
}))

vi.mock('../utils/thresholds', () => ({
  getThresholds: vi.fn(),
}))

describe('Devices - Connect Button Feature', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(devicesApi.listDevices).mockResolvedValue([])
    vi.mocked(devicesApi.getLatestTelemetry).mockResolvedValue([])
    vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({})
  })

  describe('Button Rendering Based on Mode', () => {
    it('should show Simular button when simulation mode is enabled', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': true })

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const simulateButton = wrapper.find('[title="Gerar leitura simulada"]')
      expect(simulateButton.exists()).toBe(true)
      expect(simulateButton.text()).toContain('Simular')
    })

    it('should show Conectar button when API mode is enabled (simulation disabled)', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': false })

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const connectButton = wrapper.find('[title="Conectar ao dispositivo e testar credenciais"]')
      expect(connectButton.exists()).toBe(true)
      expect(connectButton.text()).toContain('Conectar')
    })

    it('should not show both buttons at the same time', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': false })

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const simulateButton = wrapper.find('[title="Gerar leitura simulada"]')
      const connectButton = wrapper.find('[title="Conectar ao dispositivo e testar credenciais"]')

      expect(simulateButton.exists()).toBe(false)
      expect(connectButton.exists()).toBe(true)
    })
  })

  describe('Connect Functionality', () => {
    it('should call testDeviceConnection when connect button is clicked', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': false })
      vi.mocked(devicesApi.testDeviceConnection).mockResolvedValue({ power: 150.5 })

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const connectButton = wrapper.find('[title="Conectar ao dispositivo e testar credenciais"]')
      await connectButton.trigger('click')

      expect(devicesApi.testDeviceConnection).toHaveBeenCalledWith(1)
    })

    it('should update device status to online after successful connection', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': false })
      vi.mocked(devicesApi.testDeviceConnection).mockResolvedValue({ power: 150.5 })

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const connectButton = wrapper.find('[title="Conectar ao dispositivo e testar credenciais"]')
      await connectButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Check that device status was updated (in the component's data)
      const devices = (wrapper.vm as any).devices
      expect(devices[0].status).toBe('online')
    })

    it('should show loading indicator while connecting', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': false })
      
      // Make the API call hang
      let resolveConnection: (value: any) => void
      const connectionPromise = new Promise(resolve => {
        resolveConnection = resolve
      })
      vi.mocked(devicesApi.testDeviceConnection).mockReturnValue(connectionPromise as any)

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const connectButton = wrapper.find('[title="Conectar ao dispositivo e testar credenciais"]')
      await connectButton.trigger('click')
      await wrapper.vm.$nextTick()

      // Check that loading indicator is shown
      expect(connectButton.text()).toContain('⏳')
      
      // Resolve the connection
      resolveConnection!({ power: 150.5 })
      await wrapper.vm.$nextTick()
    })

    it('should display error message when connection fails', async () => {
      const mockDevice = {
        id: 1,
        user_id: 1,
        name: 'Test Device',
        status: 'offline' as const,
      }

      vi.mocked(devicesApi.listDevices).mockResolvedValue([mockDevice])
      vi.mocked(simulationMode.getAllSimulationModes).mockReturnValue({ '1': false })
      vi.mocked(devicesApi.testDeviceConnection).mockRejectedValue({
        response: {
          data: {
            error: 'device not configured for live reads'
          }
        }
      })

      const wrapper = mount(Devices, {
        global: {
          stubs: ['LoadingSpinner'],
        },
      })

      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const connectButton = wrapper.find('[title="Conectar ao dispositivo e testar credenciais"]')
      await connectButton.trigger('click')
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      // Check that error is displayed
      const error = (wrapper.vm as any).error
      expect(error).toContain('device not configured for live reads')
    })
  })
})
